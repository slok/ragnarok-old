package injection

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/log"
)

// Failer will implement the way of making a failure of a kind on a system (high level error).
type Failer interface {
	// Inject applies the failure to the system
	Inject() error

	// Revert will disable the failure
	Revert() error
}

// Injection is a failure that can be applied.
type Injection struct {
	*v1.Failure
	attacks []attack.Attacker
	ctx     context.Context
	ctxC    context.CancelFunc
	sync.Mutex
	log   log.Logger
	clock clock.Clock

	erroredAtts []attack.Attacker // Used to track the failured attacks
	appliedAtts []attack.Attacker // Used to track the correct applied attacks
}

// NewInjection Creates a new Failer object from a failure definition
// using the base global registry.
func NewInjection(f *v1.Failure, l log.Logger, cl clock.Clock) (*Injection, error) {
	return NewInjectionFromReg(f, attack.BaseReg(), l, cl)
}

// NewInjectionFromReg Creates a new SystemFailure object from a failure definition
// and a custom registry.
func NewInjectionFromReg(f *v1.Failure, reg attack.Registry, l log.Logger, cl clock.Clock) (*Injection, error) {
	// Set global logger if no logger
	if l == nil {
		l = log.Base()
	}

	if cl == nil {
		cl = clock.New()
	}

	// Create the attacks.
	atts := make([]attack.Attacker, len(f.Spec.Attacks))

	for i, tAC := range f.Spec.Attacks {
		// Check on each attack slice there is only one attack map.
		if len(tAC) != 1 {
			return nil, errors.New("configuration attack doesn't have the correct length")
		}
		// Get key/value iterating over a one element map, ugh... .
		var kind string
		var aC attack.Opts
		for kind, aC = range tAC {
			break
		}
		a, err := reg.New(kind, aC)
		if err != nil {
			return nil, err
		}
		atts[i] = a
	}

	// Update the v1.
	f.Status.Creation = cl.Now().UTC()
	f.Status.CurrentState = v1.EnabledFailureState

	ij := &Injection{
		Failure: f,
		attacks: atts,
		ctx:     context.Background(),
		log:     l,
		clock:   cl,
	}

	return ij, nil
}

// Fail implements Failer interface. Locked operation
func (i *Injection) Fail() error {
	// Set correct state and only allow execution of not executed failures
	i.Lock()
	if i.Status.CurrentState != v1.EnabledFailureState {
		return fmt.Errorf("invalid state. The only valid state for execution is: %s", v1.EnabledFailureState)
	}
	i.Status.CurrentState = v1.ExecutingFailureState
	i.Status.Executed = i.clock.Now().UTC()
	i.Unlock()

	i.ctx, i.ctxC = context.WithCancel(i.ctx)

	// channels for the attack results
	errCh := make(chan attack.Attacker)
	applyCh := make(chan attack.Attacker)

	for _, a := range i.attacks {
		go func(a attack.Attacker) {
			if err := a.Apply(i.ctx); err != nil {
				// Process the error, if there is any error then we need to revert
				log.Errorf("error aplying attack: %s", err)
				errCh <- a
			} else {
				applyCh <- a
			}
		}(a)
	}

	// Check for errors, sync with channels
	for j := 0; j < len(i.attacks); j++ {
		select {
		case a := <-errCh:
			i.erroredAtts = append(i.erroredAtts, a)
		case a := <-applyCh:
			i.appliedAtts = append(i.appliedAtts, a)
		}
	}

	// Check if there are any errors, if there are errors then revert the applied ones
	if len(i.erroredAtts) > 0 {
		// If reverting correct applied ones then return different error
		if err := i.Revert(); err != nil {
			log.Error(err)
			return fmt.Errorf("error aplying failure & error when trying to revert the applied ones")
		}
		i.Lock()
		i.Status.CurrentState = v1.ErroredFailureState
		i.Unlock()
		return fmt.Errorf("error aplying failure")
	}

	// Set execution timer and start the countdown until the revert
	go func() {
		select {
		case <-i.ctx.Done():
			i.log.Info("context on system failure done")
		case <-i.clock.After(i.Spec.Timeout):
			i.log.Info("system failure finished")
		}
		i.Lock()
		// Don't revert if not executing
		if i.Status.CurrentState != v1.ExecutingFailureState {
			i.log.Warnf("system failure attempt to finish but this is not in running state: %s", i.Status.CurrentState)
			return
		}
		i.Unlock()
		i.Revert()
	}()
	i.Status.Executed = i.clock.Now().UTC()
	i.log.Infof("execution of '%s' failure started", i.Metadata.ID)
	return nil
}

// Revert implements Revert interface.
func (i *Injection) Revert() error {
	i.log.Infof("reverting '%s' failure", i.Metadata.ID)
	defer i.ctxC()

	// Only revert the applied attacks
	errsCh := make(chan error)
	for _, a := range i.appliedAtts {
		// TODO: Retry
		go func(a attack.Attacker) {
			errsCh <- a.Revert()
		}(a)
	}

	i.Status.CurrentState = v1.DisabledFailureState
	errStr := ""
	for j := 0; j < len(i.appliedAtts); j++ {
		if err := <-errsCh; err != nil {
			errStr = fmt.Sprintf("%s; %s", errStr, err)
		}
	}

	var err error
	i.Lock()
	i.Status.Finished = i.clock.Now().UTC()
	if errStr != "" {
		i.Status.CurrentState = v1.ErroredRevertingFailureState
		err = fmt.Errorf("error reverting failure (triggered by errored attacks when aplying attacks): %s", errStr)
	}
	i.Unlock()
	return err
}
