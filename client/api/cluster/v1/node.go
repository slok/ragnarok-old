package v1

import (
	"fmt"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	apiutil "github.com/slok/ragnarok/api/util"
	"github.com/slok/ragnarok/apimachinery/validator"
	"github.com/slok/ragnarok/apimachinery/watch"
	"github.com/slok/ragnarok/client/repository"
)

var objType = api.TypeMeta{Kind: clusterv1.NodeKind, Version: clusterv1.NodeVersion}

// NodeClientInterface has the required logic to manage Nodes.
type NodeClientInterface interface {
	Create(node *clusterv1.Node) (*clusterv1.Node, error)
	Update(node *clusterv1.Node) (*clusterv1.Node, error)
	Delete(id string) error
	Get(id string) (*clusterv1.Node, error)
	List(opts api.ListOptions) (*clusterv1.NodeList, error)
	Watch(opts api.ListOptions) (watch.Watcher, error)
	// TODO Patch
}

// NodeClient has the required logic to manage Nodes.
type NodeClient struct {
	validator validator.ObjectValidator
	repoCli   repository.Client
}

// NewNodeClient returns a new NodeClient.
func NewNodeClient(validator validator.ObjectValidator, repoCli repository.Client) *NodeClient {
	return &NodeClient{
		validator: validator,
		repoCli:   repoCli,
	}
}

func (n *NodeClient) typeAssertNode(obj api.Object) (*clusterv1.Node, error) {
	node, ok := obj.(*clusterv1.Node)
	if !ok {
		return nil, fmt.Errorf("could not make the type assertion from obj to node. Wrong type")
	}
	return node, nil
}
func (n *NodeClient) typeAssertNodeList(objs api.ObjectList) (*clusterv1.NodeList, error) {
	nodes := make([]*clusterv1.Node, len(objs.GetItems()))
	for i, obj := range objs.GetItems() {
		node, ok := obj.(*clusterv1.Node)
		if !ok {
			return nil, fmt.Errorf("could not make the type assertion from obj to node. Wrong type")
		}
		nodes[i] = node
	}
	nList := clusterv1.NewNodeList(nodes, objs.GetListMetadata().Continue)
	return &nList, nil
}

func (n *NodeClient) validate(node *clusterv1.Node) error {
	// Check valid object.
	if errs := n.validator.Validate(node); len(errs) > 0 {
		return fmt.Errorf("error on validation: %s", errs)
	}
	return nil
}

// Create satisfies NodeClientInterface interface.
func (n *NodeClient) Create(node *clusterv1.Node) (*clusterv1.Node, error) {
	// Check valid object.
	if err := n.validate(node); err != nil {
		return nil, err
	}

	obj, err := n.repoCli.Create(node)
	if err != nil {
		return nil, err
	}
	return n.typeAssertNode(obj)
}

// Update satisfies NodeClientInterface interface.
func (n *NodeClient) Update(node *clusterv1.Node) (*clusterv1.Node, error) {
	// Check valid object.
	if err := n.validate(node); err != nil {
		return nil, err
	}

	obj, err := n.repoCli.Update(node)
	if err != nil {
		return nil, err
	}
	return n.typeAssertNode(obj)
}

// Delete satisfies NodeClientInterface interface.
func (n *NodeClient) Delete(id string) error {
	// get the full ID
	fullID := apiutil.GetFullIDFromType(objType, id)
	return n.repoCli.Delete(fullID)
}

// Get satisfies NodeClientInterface interface.
func (n *NodeClient) Get(id string) (*clusterv1.Node, error) {
	fullID := apiutil.GetFullIDFromType(objType, id)
	obj, err := n.repoCli.Get(fullID)
	if err != nil {
		return nil, err
	}
	return n.typeAssertNode(obj)
}

// List satisfies NodeClientInterface interface.
func (n *NodeClient) List(opts api.ListOptions) (*clusterv1.NodeList, error) {
	opts.TypeMeta = clusterv1.NodeTypeMeta
	objs, err := n.repoCli.List(opts)
	if err != nil {
		return nil, err
	}
	return n.typeAssertNodeList(objs)
}

// Watch satisfies NodeClientInterface interface.
func (n *NodeClient) Watch(opts api.ListOptions) (watch.Watcher, error) {
	opts.TypeMeta = clusterv1.NodeTypeMeta
	return n.repoCli.Watch(opts)
}
