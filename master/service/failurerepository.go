package service

import (
	"sync"

	"github.com/slok/ragnarok/master/model"
)

// FailureRepository is the way the master keeps track of the failures.
type FailureRepository interface {
	// Store adds a failure to the registry.
	Store(failure *model.Failure) error

	// Delete deletes a failure from the registry.
	Delete(id string)

	// Get gets a failure from the registry.
	Get(id string) (*model.Failure, bool)

	// GetAll gets all the failures from the registry.
	GetAll() map[string]*model.Failure
}

// MemFailureRepository is a represententation of the failure regsitry using a memory map.
type MemFailureRepository struct {
	reg map[string]*model.Failure
	sync.Mutex
}

// NewMemFailureRepository returns a new MemFailureRepository
func NewMemFailureRepository() *MemFailureRepository {
	return &MemFailureRepository{
		reg: map[string]*model.Failure{},
	}
}

// Store satisfies FailureRepository interface.
func (m *MemFailureRepository) Store(failure *model.Failure) error {
	m.Lock()
	defer m.Unlock()
	m.reg[failure.ID] = failure

	return nil
}

// Delete satisfies FailureRepository interface.
func (m *MemFailureRepository) Delete(id string) {
	m.Lock()
	defer m.Unlock()

	delete(m.reg, id)
}

// Get satisfies FailureRepository interface.
func (m *MemFailureRepository) Get(id string) (*model.Failure, bool) {
	m.Lock()
	defer m.Unlock()

	f, ok := m.reg[id]

	return f, ok
}

// GetAll satisfies FailureRepository interface.
func (m *MemFailureRepository) GetAll() map[string]*model.Failure {
	m.Lock()
	defer m.Unlock()

	return m.reg
}
