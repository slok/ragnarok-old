package master

import (
	"sync"
)

// Node is an internal and simplified representation of a failure node on the masters
// TODO: Rethink the reuse of node
type Node struct {
	ID   string            // ID is the id of the node
	Tags map[string]string // Tags are the tags related with the node
}

// NodeRepository is the way the master should store the nodes
type NodeRepository interface {
	// StoreNode adds a node to the registry
	StoreNode(id string, node *Node) error

	// DeleteNode deletes a node from the registry
	DeleteNode(id string)

	// GetNode gets a node from the registry
	GetNode(id string) (*Node, bool)

	// GetNodes gets all the nodes from the registry
	GetNodes() map[string]*Node
}

// NewMemNodeRepository returns a new memory node registry
func NewMemNodeRepository() *MemNodeRepository {
	return &MemNodeRepository{
		reg: map[string]*Node{},
	}
}

// MemNodeRepository is a representation of the node registry using memorymap
type MemNodeRepository struct {
	reg map[string]*Node
	sync.Mutex
}

// StoreNode satisfies NodeRepository interface
func (m *MemNodeRepository) StoreNode(id string, node *Node) error {
	m.Lock()
	defer m.Unlock()

	m.reg[id] = node

	return nil
}

// DeleteNode satisfies NodeRepository interface
func (m *MemNodeRepository) DeleteNode(id string) {
	m.Lock()
	defer m.Unlock()

	delete(m.reg, id)
}

// GetNode satisfies GetNode interface
func (m *MemNodeRepository) GetNode(id string) (*Node, bool) {
	m.Lock()
	defer m.Unlock()
	n, ok := m.reg[id]
	return n, ok
}

// GetNodes satisfies GetNode interface
func (m *MemNodeRepository) GetNodes() map[string]*Node {
	m.Lock()
	defer m.Unlock()
	return m.reg
}
