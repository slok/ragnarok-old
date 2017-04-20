package attack

import (
	"errors"
	"fmt"
)

var registry = make(map[string]Creator)

// Creator is a custom type used for attack creation.
type Creator func(opts map[string]interface{}) (Attacker, error)

// Register registers a attack creator.
func Register(id string, c Creator) error {
	if id == "" {
		return errors.New("invalid id for attacker registration")
	}
	registry[id] = c
	return nil
}

// Deregister deregisters a attack creator.
func Deregister(id string) error {
	if !Exists(id) {
		return fmt.Errorf("'%s' is and invalid id for attacker deregistration", id)
	}
	delete(registry, id)
	return nil
}

// Exists returns true if the creator is registered for a given ID, false instead.
func Exists(id string) bool {
	_, ok := registry[id]
	return ok
}

// New is the factory of the attacks based on IDs and options.
func New(id string, opts map[string]interface{}) (Attacker, error) {
	c, ok := registry[id]
	if !ok {
		return nil, fmt.Errorf("%s is not a correct Attack", id)
	}
	return c(opts)
}
