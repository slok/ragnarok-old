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
	GetAll() []*model.Failure

	// GetAllByNode gets all the failures of a node from the registry.
	GetAllByNode(nodeID string) []*model.Failure
}

// MemFailureRepository is a represententation of the failure regsitry using a memory map.
type MemFailureRepository struct {
	reg       map[string]*model.Failure
	regByNode map[string]map[string]*model.Failure
	sync.Mutex
}

// NewMemFailureRepository returns a new MemFailureRepository
func NewMemFailureRepository() *MemFailureRepository {
	return &MemFailureRepository{
		reg:       map[string]*model.Failure{},
		regByNode: map[string]map[string]*model.Failure{},
	}
}

// Store satisfies FailureRepository interface.
func (m *MemFailureRepository) Store(failure *model.Failure) error {
	m.Lock()
	defer m.Unlock()
	m.reg[failure.ID] = failure
	if _, ok := m.regByNode[failure.NodeID]; !ok {
		m.regByNode[failure.NodeID] = map[string]*model.Failure{}
	}
	m.regByNode[failure.NodeID][failure.ID] = failure

	return nil
}

// Delete satisfies FailureRepository interface.
func (m *MemFailureRepository) Delete(id string) {
	m.Lock()
	defer m.Unlock()

	f, ok := m.reg[id]
	if !ok {
		return
	}

	delete(m.reg, id)
	delete(m.regByNode[f.NodeID], id)
}

// Get satisfies FailureRepository interface.
func (m *MemFailureRepository) Get(id string) (*model.Failure, bool) {
	m.Lock()
	defer m.Unlock()

	f, ok := m.reg[id]

	return f, ok
}

// GetAll satisfies FailureRepository interface.
func (m *MemFailureRepository) GetAll() []*model.Failure {
	m.Lock()
	defer m.Unlock()
	res := []*model.Failure{}
	for _, f := range m.reg {
		res = append(res, f)
	}
	return res
}

// GetAllByNode satisfies FailureRepository interface.
func (m *MemFailureRepository) GetAllByNode(nodeID string) []*model.Failure {
	m.Lock()
	defer m.Unlock()
	res := []*model.Failure{}
	tmpReg, ok := m.regByNode[nodeID]
	if ok {
		for _, f := range tmpReg {
			res = append(res, f)
		}
	}
	return res
}
