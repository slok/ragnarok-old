package service

import (
	"sync"

	"github.com/slok/ragnarok/failure"
)

// FailureRepository is the way the master keeps track of the failures.
type FailureRepository interface {
	// Store adds a failure to the registry.
	Store(failure *failure.Failure) error

	// Delete deletes a failure from the registry.
	Delete(id string)

	// Get gets a failure from the registry.
	Get(id string) (*failure.Failure, bool)

	// GetAll gets all the failures from the registry.
	GetAll() []*failure.Failure

	// GetAllByNode gets all the failures of a node from the registry.
	GetAllByNode(nodeID string) []*failure.Failure
}

// MemFailureRepository is a represententation of the failure regsitry using a memory map.
type MemFailureRepository struct {
	reg       map[string]*failure.Failure
	regByNode map[string]map[string]*failure.Failure
	sync.Mutex
}

// NewMemFailureRepository returns a new MemFailureRepository
func NewMemFailureRepository() *MemFailureRepository {
	return &MemFailureRepository{
		reg:       map[string]*failure.Failure{},
		regByNode: map[string]map[string]*failure.Failure{},
	}
}

// Store satisfies FailureRepository interface.
func (m *MemFailureRepository) Store(f *failure.Failure) error {
	m.Lock()
	defer m.Unlock()
	m.reg[f.ID] = f
	if _, ok := m.regByNode[f.NodeID]; !ok {
		m.regByNode[f.NodeID] = map[string]*failure.Failure{}
	}
	m.regByNode[f.NodeID][f.ID] = f

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
func (m *MemFailureRepository) Get(id string) (*failure.Failure, bool) {
	m.Lock()
	defer m.Unlock()

	f, ok := m.reg[id]

	return f, ok
}

// GetAll satisfies FailureRepository interface.
func (m *MemFailureRepository) GetAll() []*failure.Failure {
	m.Lock()
	defer m.Unlock()
	res := []*failure.Failure{}
	for _, f := range m.reg {
		res = append(res, f)
	}
	return res
}

// GetAllByNode satisfies FailureRepository interface.
func (m *MemFailureRepository) GetAllByNode(nodeID string) []*failure.Failure {
	m.Lock()
	defer m.Unlock()
	res := []*failure.Failure{}
	tmpReg, ok := m.regByNode[nodeID]
	if ok {
		for _, f := range tmpReg {
			res = append(res, f)
		}
	}
	return res
}
