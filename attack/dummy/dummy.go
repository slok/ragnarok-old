package dummy

import (
	"context"

	"github.com/slok/ragnarok/attack"
)

const (
	// DummyID is the identifier of the attack
	DummyID = "dummy"
)

// Register the creator of the attack
func init() {
	attack.Register(DummyID, attack.CreatorFunc(func(o attack.Opts) (attack.Attacker, error) {
		return NewDummy(o)
	}))
}

// Dummy failer will do nothing.
type Dummy struct{}

// NewDummy returns a new Dummy attack.
func NewDummy(_ attack.Opts) (*Dummy, error) {
	return &Dummy{}, nil
}

// Apply will do nothing.
func (d *Dummy) Apply(_ context.Context) error { return nil }

// Revert will do nothing.
func (d *Dummy) Revert() error { return nil }
