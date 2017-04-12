package attack

import "context"

// Attacker implements the lowest level of error applied to a system
type Attacker interface {

	// Apply applies an attack or fault to the system
	Apply(ctx context.Context) error

	// Revert reverts an attack or fault from the system
	Revert() error
}
