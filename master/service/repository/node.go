package repository

import (
	"sync"

	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
)

// Node is the way the master should store the nodes.
type Node interface {
	// StoreNode adds a node to the registry.
	StoreNode(id string, node clusterv1.Node) error

	// DeleteNode deletes a node from the registry.
	DeleteNode(id string)

	// GetNode gets a node from the registry.
	GetNode(id string) (*clusterv1.Node, bool)

	// GetNodes gets all the nodes from the registry.
	GetNodes() map[string]*clusterv1.Node

	// GetNodesByLabels gets all the nodes from the registry using labels.
	GetNodesByLabels(labels clusterv1.NodeLabels) map[string]*clusterv1.Node
}

// NewMemNode returns a new memory node registry.
func NewMemNode() *MemNode {
	return &MemNode{
		reg: map[string]*clusterv1.Node{},
	}
}

// MemNode is a representation of the node registry using memorymap, used only
// as a first implementation to get working the first version.
type MemNode struct {
	reg map[string]*clusterv1.Node
	sync.Mutex
}

// StoreNode satisfies Node interface.
func (m *MemNode) StoreNode(id string, node clusterv1.Node) error {
	m.Lock()
	defer m.Unlock()

	m.reg[id] = &node

	return nil
}

// DeleteNode satisfies Node interface.
func (m *MemNode) DeleteNode(id string) {
	m.Lock()
	defer m.Unlock()

	delete(m.reg, id)
}

// GetNode satisfies Node interface.
func (m *MemNode) GetNode(id string) (*clusterv1.Node, bool) {
	m.Lock()
	defer m.Unlock()
	n, ok := m.reg[id]
	return n, ok
}

// GetNodes satisfies Node interface.
func (m *MemNode) GetNodes() map[string]*clusterv1.Node {
	m.Lock()
	defer m.Unlock()
	return m.reg
}

// GetNodesByLabels satisfies Node interface.
func (m *MemNode) GetNodesByLabels(labels clusterv1.NodeLabels) map[string]*clusterv1.Node {
	result := map[string]*clusterv1.Node{}
	if len(labels) == 0 {
		return result
	}

	m.Lock()
	defer m.Unlock()

NodeLoop:
	for nName, node := range m.reg {
		// Check on each node if staisfies the labels.
		for lk, lv := range labels {
			if nv, ok := node.Spec.Labels[lk]; !ok || (ok && nv != lv) {
				continue NodeLoop // Continue next node iteration, not valid.
			}
		}
		result[nName] = node
	}
	return result
}
