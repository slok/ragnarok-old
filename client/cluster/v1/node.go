package v1

import (
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/apimachinery/watch"
)

// NodeListOptions are the options required to list nodes
type NodeListOptions struct {
	Selector map[string]string
}

// Node has the required logic to manage Nodes.
type Node interface {
	Create(node *clusterv1.Node) (*clusterv1.Node, error)
	Update(node *clusterv1.Node) (*clusterv1.Node, error)
	Delete(id string) error
	Get(id string) (*clusterv1.Node, error)
	List(opts NodeListOptions) ([]*clusterv1.Node, error)
	Watch(opts NodeListOptions) (watch.Watch, error)
	// TODO Patch
}
