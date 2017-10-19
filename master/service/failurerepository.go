package service

import (
	"sync"

	"github.com/slok/ragnarok/api/chaos/v1"
)

// FailureRepository is the way the master keeps track of the failures.
type FailureRepository interface {
	// Store adds a failure to the registry.
	Store(failure *v1.Failure) error

	// Delete deletes a failure from the registry.
	Delete(id string)

	// Get gets a failure from the registry.
	Get(id string) (*v1.Failure, bool)

	// GetAll gets all the failures from the registry.
	GetAll() []*v1.Failure

	// GetAllByNode gets all the failures of a node from the registry.
	GetAllByNode(nodeID string) []*v1.Failure

	// GetNotStaleByNode gets all not stale failures of a node from the registry.
	GetNotStaleByNode(nodeID string) []*v1.Failure
}

// MemFailureRepository is a represententation of the failure regsitry using a memory map.
type MemFailureRepository struct {
	reg       map[string]*v1.Failure
	regByNode map[string]map[string]*v1.Failure
	sync.Mutex
}

// NewMemFailureRepository returns a new MemFailureRepository
func NewMemFailureRepository() *MemFailureRepository {
	return &MemFailureRepository{
		reg:       map[string]*v1.Failure{},
		regByNode: map[string]map[string]*v1.Failure{},
	}
}

// Store satisfies FailureRepository interface.
func (m *MemFailureRepository) Store(f *v1.Failure) error {
	m.Lock()
	defer m.Unlock()
	m.reg[f.Metadata.ID] = f
	if _, ok := m.regByNode[f.Metadata.NodeID]; !ok {
		m.regByNode[f.Metadata.NodeID] = map[string]*v1.Failure{}
	}
	m.regByNode[f.Metadata.NodeID][f.Metadata.ID] = f

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
	delete(m.regByNode[f.Metadata.NodeID], id)
}

// Get satisfies FailureRepository interface.
func (m *MemFailureRepository) Get(id string) (*v1.Failure, bool) {
	m.Lock()
	defer m.Unlock()

	f, ok := m.reg[id]

	return f, ok
}

// GetAll satisfies FailureRepository interface.
func (m *MemFailureRepository) GetAll() []*v1.Failure {
	m.Lock()
	defer m.Unlock()
	res := []*v1.Failure{}
	for _, f := range m.reg {
		res = append(res, f)
	}
	return res
}

// getAllByNode gets all the failures with an stale filter, if the filter is true
// then it will return also the stale ones, if not then it will return all expect the stale ones.
func (m *MemFailureRepository) getAllByNode(nodeID string, stale bool) []*v1.Failure {
	m.Lock()
	defer m.Unlock()
	res := []*v1.Failure{}
	tmpReg, ok := m.regByNode[nodeID]
	if ok {
		for _, f := range tmpReg {
			// Only add the ones that we want, do we want stale data?
			if !stale && f.Status.CurrentState == v1.StaleFailureState {
				continue
			}
			res = append(res, f)
		}
	}
	return res
}

// GetAllByNode satisfies FailureRepository interface.
func (m *MemFailureRepository) GetAllByNode(nodeID string) []*v1.Failure {
	return m.getAllByNode(nodeID, true)
}

// GetNotStaleByNode satisfies FailureRepository interface.
func (m *MemFailureRepository) GetNotStaleByNode(nodeID string) []*v1.Failure {
	return m.getAllByNode(nodeID, false)
}
