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

// NodeRegistry is the way the master should store the nodes
type NodeRegistry interface {
	// AddNode adds a node to the registry
	AddNode(id string, node *Node) error

	// DeleteNode deletes a node from the registry
	DeleteNode(id string)

	// GetNode gets a node from the registry
	GetNode(id string) (*Node, bool)

	// GetNodes gets all the nodes from the registry
	GetNodes() map[string]*Node
}

// NewMemNodeRegistry returns a new memory node registry
func NewMemNodeRegistry() *MemNodeRegistry {
	return &MemNodeRegistry{
		reg: map[string]*Node{},
	}
}

// MemNodeRegistry is a representation of the node registry using memorymap
type MemNodeRegistry struct {
	reg map[string]*Node
	sync.Mutex
}

// AddNode satisfies NodeRegistry interface
func (m *MemNodeRegistry) AddNode(id string, node *Node) error {
	m.Lock()
	defer m.Unlock()

	m.reg[id] = node

	return nil
}

// DeleteNode satisfies NodeRegistry interface
func (m *MemNodeRegistry) DeleteNode(id string) {
	m.Lock()
	defer m.Unlock()

	delete(m.reg, id)
}

// GetNode satisfies GetNode interface
func (m *MemNodeRegistry) GetNode(id string) (*Node, bool) {
	m.Lock()
	defer m.Unlock()
	n, ok := m.reg[id]
	return n, ok
}

// GetNodes satisfies GetNode interface
func (m *MemNodeRegistry) GetNodes() map[string]*Node {
	m.Lock()
	defer m.Unlock()
	return m.reg
}
