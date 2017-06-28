package failure

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/log"
)

// State is the current situation of the failure.
type State int

const (
	// Created means a not executed failure.
	Created State = iota
	// Executing means an executing failure.
	Executing
	// Reverted means a reverted failure.
	Reverted
	// Error means a failure that errored.
	Error
	// ErrorReverting means a failure that errored reverting it.
	ErrorReverting
	// Unknown means a staet of the failure that is unknown.
	Unknown
)

func (s State) String() string {
	switch s {
	case Created:
		return "created"
	case Executing:
		return "executing"
	case Reverted:
		return "reverted"
	case Error:
		return "error"
	case ErrorReverting:
		return "error reverting"
	default:
		return "unknown"
	}
}

// Failer will implement the way of making a failure of a kind on a system (high level error).
type Failer interface {
	// Fail applies the failure to the system
	Fail() error

	// Revert will disable the failure.
	Revert() error
}

// SystemFailure is the most basic kind of failure.
type SystemFailure struct {
	id       string
	timeout  time.Duration
	attacks  []attack.Attacker
	ctx      context.Context
	ctxC     context.CancelFunc
	creation time.Time
	executed time.Time
	finished time.Time
	State    State
	sync.Mutex
	log   log.Logger
	clock clock.Clock

	erroredAtts []attack.Attacker // Used to track the failured attacks
	appliedAtts []attack.Attacker // Used to track the correct applied attacks
}

// NewSystemFailure Creates a new SystemFailure object from a failure definition
// using the base global registry.
func NewSystemFailure(c Config, l log.Logger, cl clock.Clock) (*SystemFailure, error) {
	return NewSystemFailureFromReg(c, attack.BaseReg(), l, cl)
}

// NewSystemFailureFromReg Creates a new SystemFailure object from a failure definition
// and a custom registry.
func NewSystemFailureFromReg(c Config, reg attack.Registry, l log.Logger, cl clock.Clock) (*SystemFailure, error) {
	// Set global logger if no logger
	if l == nil {
		l = log.Base()
	}

	if cl == nil {
		cl = clock.New()
	}

	// Create the attacks.
	atts := make([]attack.Attacker, len(c.Attacks))

	for i, tAC := range c.Attacks {
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

	id := "random_id" // TODO
	f := &SystemFailure{
		id:       id,
		timeout:  c.Timeout,
		attacks:  atts,
		creation: cl.Now().UTC(),
		ctx:      context.Background(),
		State:    Created,
		log:      l,
		clock:    cl,
	}

	return f, nil
}

// Fail implements Failer interface. Locked operation
func (s *SystemFailure) Fail() error {
	// Set correct state and only allow execution of not executed failures
	s.Lock()
	if s.State != Created {
		return fmt.Errorf("invalid state. The only valid state for execution is: %s", Created)
	}
	s.State = Executing
	s.executed = s.clock.Now().UTC()
	s.Unlock()

	s.ctx, s.ctxC = context.WithCancel(s.ctx)

	// channels for the attack results
	errCh := make(chan attack.Attacker)
	applyCh := make(chan attack.Attacker)

	for _, a := range s.attacks {
		go func(a attack.Attacker) {
			if err := a.Apply(s.ctx); err != nil {
				// Process the error, if there is any error then we need to revert
				log.Errorf("error aplying attack: %s", err)
				errCh <- a
			} else {
				applyCh <- a
			}
		}(a)
	}

	// Check for errors, sync with channels
	for i := 0; i < len(s.attacks); i++ {
		select {
		case a := <-errCh:
			s.erroredAtts = append(s.erroredAtts, a)
		case a := <-applyCh:
			s.appliedAtts = append(s.appliedAtts, a)
		}
	}

	// Check if there are any errors, if there are errors then revert the applied ones
	if len(s.erroredAtts) > 0 {
		// If reverting correct applied ones then return different error
		if err := s.Revert(); err != nil {
			log.Error(err)
			return fmt.Errorf("error aplying failure & error when trying to revert the applied ones")
		}
		s.Lock()
		s.State = Error
		s.Unlock()
		return fmt.Errorf("error aplying failure")
	}

	// Set execution timer and start the countdown until the revert
	go func() {
		select {
		case <-s.ctx.Done():
			s.log.Info("Context on system failure done")
		case <-s.clock.After(s.timeout):
			s.log.Info("System failure finished")
		}
		s.Lock()
		// Don't revert if not executing
		if s.State != Executing {
			s.log.Warnf("System failure attempt to finish but this is not in running state: %s", s.State)
			return
		}
		s.Unlock()
		s.Revert()
	}()
	s.executed = s.clock.Now().UTC()
	s.log.Infof("Execution of '%s' failure started", s.id)
	return nil
}

// Revert implements Revert interface.
func (s *SystemFailure) Revert() error {
	s.log.Infof("Reverting '%s' failure", s.id)
	defer s.ctxC()

	// Only revert the applied attacks
	errsCh := make(chan error)
	for _, a := range s.appliedAtts {
		// TODO: Retry
		go func(a attack.Attacker) {
			errsCh <- a.Revert()
		}(a)
	}

	s.State = Reverted
	errStr := ""
	for i := 0; i < len(s.appliedAtts); i++ {
		if err := <-errsCh; err != nil {
			errStr = fmt.Sprintf("%s; %s", errStr, err)
		}
	}

	var err error
	s.Lock()
	s.finished = s.clock.Now().UTC()
	if errStr != "" {
		s.State = ErrorReverting
		err = fmt.Errorf("error reverting failure (triggered by errored attacks when aplying attacks): %s", errStr)
	}
	s.Unlock()
	return err
}
