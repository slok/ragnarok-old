package failure

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/slok/ragnarok/attack"
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
	// Unknown means a staet of the failure that is unknown.
	Unknown
)

// Failer will implement the way of making a failure of a kind on a system (high level error).
type Failer interface {
	// Fail applies the failure to the system for a desired time (or forever if is 0).
	Fail(duration time.Duration) error

	// Revert will disable the failure.
	Revert() error
}

// SystemFailure is the most basic kind of failure.
type SystemFailure struct {
	id       string
	timeout  time.Duration
	attacks  []attack.Attacker
	ctx      context.Context
	creation time.Time
	executed time.Time
	finished time.Time
	State    State
	sync.Mutex
	log log.Logger
}

// NewSystemFailure Creates a new SystemFailure object from a failure definition
// using the base global registry.
func NewSystemFailure(c Config, l log.Logger) (*SystemFailure, error) {
	return NewSystemFailureFromReg(c, attack.BaseReg(), l)
}

// NewSystemFailureFromReg Creates a new SystemFailure object from a failure definition
// and a custom registry.
func NewSystemFailureFromReg(c Config, reg attack.Registry, l log.Logger) (*SystemFailure, error) {
	// Set global logger if no logger
	if l == nil {
		l = log.Base()
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
		creation: time.Now().UTC(),
		ctx:      context.Background(),
		State:    Created,
	}

	return f, nil
}
