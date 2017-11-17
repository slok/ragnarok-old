package v1

import (
	"fmt"
	"sync"

	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/apimachinery/validator"
	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/client/util"
)

// NodeMem is a client implementation that stores in memory (used
// for fast development). It will have storage logic also.
type NodeMem struct {
	validator validator.ObjectValidator
	reg       map[string]*clusterv1.Node
	sync.Mutex
}

// NewDefaultNodeMem returns a default node memory repository.
func NewDefaultNodeMem() *NodeMem {
	return &NodeMem{
		validator: validator.DefaultObject,
		reg:       map[string]*clusterv1.Node{},
	}
}

// NewNodeMem returns a new node repository client
func NewNodeMem(validator validator.ObjectValidator, registry map[string]*clusterv1.Node) *NodeMem {
	return &NodeMem{
		validator: validator,
		reg:       registry,
	}
}

// Create stores a node in memory and returns an error if already exists.
func (n *NodeMem) Create(node *clusterv1.Node) (*clusterv1.Node, error) {
	// Check valid object.
	if errs := n.validator.Validate(node); len(errs) > 0 {
		return nil, fmt.Errorf("error on validation: %s", errs)
	}

	n.Lock()
	defer n.Unlock()

	if _, ok := n.reg[node.Metadata.ID]; ok {
		return nil, fmt.Errorf("node %s already present", node.Metadata.ID)
	}
	n.reg[node.Metadata.ID] = node

	return node, nil
}

// Update updates a node in memory, it will error if doesn't exists, the new obj will
// be stored and the old one will be deleted (is not a patch).
func (n *NodeMem) Update(node *clusterv1.Node) (*clusterv1.Node, error) {
	// Check valid object.
	if errs := n.validator.Validate(node); len(errs) > 0 {
		return nil, fmt.Errorf("error on validation: %s", errs)
	}

	n.Lock()
	defer n.Unlock()

	if _, ok := n.reg[node.Metadata.ID]; !ok {
		return nil, fmt.Errorf("node %s not present", node.Metadata.ID)
	}
	n.reg[node.Metadata.ID] = node

	return node, nil
}

// Delete will delete an object. It will not return an error if the object doesn't exists.
func (n *NodeMem) Delete(id string) error {
	n.Lock()
	defer n.Unlock()

	delete(n.reg, id)
	return nil
}

// Get will get an object based on its id.
func (n *NodeMem) Get(id string) (*clusterv1.Node, error) {
	n.Lock()
	defer n.Unlock()
	node, ok := n.reg[id]
	if !ok {
		return nil, fmt.Errorf("node %s not present", id)
	}
	return node, nil
}

func (n *NodeMem) listAll() ([]*clusterv1.Node, error) {
	n.Lock()
	defer n.Unlock()
	nl := []*clusterv1.Node{}
	for _, node := range n.reg {
		nl = append(nl, node)
	}

	return nl, nil
}

func (n *NodeMem) listBySelector(selector map[string]string) ([]*clusterv1.Node, error) {
	n.Lock()
	defer n.Unlock()
	nl := []*clusterv1.Node{}

	for _, node := range n.reg {
		if util.SelectorMatchesLabels(node.Metadata.Labels, selector) {
			nl = append(nl, node)
		}
	}
	return nl, nil
}

// List will return a list of objs based on the options provided.
func (n *NodeMem) List(opts NodeListOptions) ([]*clusterv1.Node, error) {
	// Return all
	if len(opts.Selector) == 0 {
		return n.listAll()
	}

	// Return filtered.
	return n.listBySelector(opts.Selector)

}

// Watch will return a channel where the new objects will be sent.
func (n *NodeMem) Watch(opts NodeListOptions) (watch.Watch, error) {
	return nil, fmt.Errorf("Not implemented")
}
