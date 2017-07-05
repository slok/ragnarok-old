package attack

import (
	"errors"
	"fmt"

	"github.com/slok/ragnarok/log"
)

// Creator interface stuff.

// Opts are the option map of an attack.
type Opts map[string]interface{}

// Creater should be implement by any attacker that wants to be automatically created by the global
// registry of attacks.
type Creater interface {
	Create(opts Opts) (Attacker, error)
}

// CreatorFunc implements Creater interface as a handy way of creating creaters quickly
type CreatorFunc func(o Opts) (Attacker, error)

// Create implements Creater.
func (c CreatorFunc) Create(opts Opts) (Attacker, error) {
	return c(opts)
}

// Registry interface stuff.

// Registry is an interface that needs to be implemented by any attack registry.
type Registry interface {
	Register(id string, c Creater) error
	Deregister(id string) error
	Exists(id string) bool
	New(id string, opts Opts) (Attacker, error)
}

// SimpleRegistry is basic registry of attackers.
type SimpleRegistry map[string]Creater

// NewSimpleRegistry returns a basic registerer.
func NewSimpleRegistry() SimpleRegistry {
	r := (SimpleRegistry)(make(map[string]Creater))
	return r
}

// Register registers a attack creator.
func (r SimpleRegistry) Register(id string, c Creater) error {
	if id == "" {
		return errors.New("invalid id for attacker registration")
	}
	r[id] = c
	log.With("attack", id).Info("attack registered")
	return nil
}

// Deregister deregisters a attack creator.
func (r SimpleRegistry) Deregister(id string) error {
	if !r.Exists(id) {
		return fmt.Errorf("'%s' is and invalid id for attacker deregistration", id)
	}
	delete(r, id)
	return nil
}

// Exists returns true if the creator is registered for a given ID, false instead.
func (r SimpleRegistry) Exists(id string) bool {
	_, ok := r[id]
	return ok
}

// New is the factory of the attacks based on IDs and options.
func (r SimpleRegistry) New(id string, opts Opts) (Attacker, error) {
	c, ok := r[id]
	if !ok {
		return nil, fmt.Errorf("%s is not a correct Attack", id)
	}
	return c.Create(opts)
}

// Global tools for the main registry.

// baseReg is the common global registry used when a custom registry is not used.
var baseReg = NewSimpleRegistry()

// BaseReg returns global application attack creator registry.
func BaseReg() Registry {
	return baseReg
}

// Register registers a new attack creator on base registry.
func Register(id string, c Creater) error {
	return baseReg.Register(id, c)
}

// Deregister deregisters a attack creator on base registry.
func Deregister(id string) error {
	return baseReg.Deregister(id)
}

// Exists checks a attack creator exists on the base registry.
func Exists(id string) bool {
	return baseReg.Exists(id)
}

// New is the factory method of attackers on the base registry.
func New(id string, opts Opts) (Attacker, error) {
	return baseReg.New(id, opts)
}
