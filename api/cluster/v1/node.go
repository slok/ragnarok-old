package v1

import (
	"time"

	"github.com/slok/ragnarok/api"
)

const (
	// NodeKind is the kind a node.
	NodeKind = "node"
	// NodeVersion is the version of the node.
	NodeVersion = "cluster/v1"

	nodeListKind    = "nodeList"
	nodeListVersion = "cluster/v1"
)

// NodeTypeMeta is the node type metadata.
var NodeTypeMeta = api.TypeMeta{
	Kind:    NodeKind,
	Version: NodeVersion,
}

// NodeListTypeMeta is the node list type metadata.
var NodeListTypeMeta = api.TypeMeta{
	Kind:    nodeListKind,
	Version: nodeListVersion,
}

// NodeState is the reprensetation of the node state.
type NodeState int

const (
	// UnknownNodeState is an unknown status.
	UnknownNodeState NodeState = iota
	// ReadyNodeState is the state when a node is ready for accepting attacks.
	ReadyNodeState
	// AttackingNodeState is the state when a node is aplying an attack.
	AttackingNodeState
	// RevertingNodeState is the state when a node is reverting an applied attack.
	RevertingNodeState
	// ErroredNodeState is the state when a node is in error state.
	ErroredNodeState
)

// String implements the stringer interface.
func (n NodeState) String() string {
	switch n {
	case ReadyNodeState:
		return "ready"
	case AttackingNodeState:
		return "attacking"
	case RevertingNodeState:
		return "reverting"
	case ErroredNodeState:
		return "errored"
	default:
		return "unknown"
	}
}

// NodeSpec has the node specific fields.
type NodeSpec struct{}

// NodeStatus has the state fo the node.
type NodeStatus struct {
	State    NodeState `json:"state,omitempty"`    // State is the state of the Node
	Creation time.Time `json:"creation,omitempty"` // Creation is when the creation of the node happenned.
}

// NewNode is a plain Node object contructor.
func NewNode() Node {
	return Node{
		TypeMeta: api.TypeMeta{
			Kind:    NodeKind,
			Version: NodeVersion,
		},
	}
}

// Node is an internal and simplified representation of a failure node on the masters
type Node struct {
	api.TypeMeta `json:",inline"`

	Metadata api.ObjectMeta `json:"metadata,omitempty"`
	Spec     NodeSpec       `json:"spec,omitempty"`
	Status   NodeStatus     `json:"status,omitempty"`
}

// GetObjectMetadata satisfies object interface.
func (n *Node) GetObjectMetadata() api.ObjectMeta {
	return n.Metadata
}

// DeepCopy satisfies object interface.
func (n *Node) DeepCopy() api.Object {
	copy := *n
	return &copy
}

// NodeList is a node list.
type NodeList struct {
	api.TypeMeta `json:",inline"`
	ListMetadata api.ListMeta `json:"listMetadata,omitempty"`
	Items        []*Node      `json:"items,omitempty"`
}

// NewNodeList returns a new NodeList.
func NewNodeList(nodes []*Node, continueList string) NodeList {
	return NodeList{
		TypeMeta: NodeListTypeMeta,
		ListMetadata: api.ListMeta{
			Continue: continueList,
		},
		Items: nodes,
	}
}

// GetObjectMetadata satisfies object interface.
func (n *NodeList) GetObjectMetadata() api.ObjectMeta {
	return api.NoObjectMeta
}

// GetListMetadata satisfies objectList interface.
func (n *NodeList) GetListMetadata() api.ListMeta {
	return n.ListMetadata
}

// GetItems satisfies ObjectList interface.
func (n *NodeList) GetItems() []api.Object {
	res := make([]api.Object, len(n.Items))
	for i, item := range n.Items {
		res[i] = api.Object(item)
	}
	return res
}

// DeepCopy satisfies object interface.
func (n *NodeList) DeepCopy() api.Object {
	ns := []*Node{}
	for i, node := range n.Items {
		n := *node
		ns[i] = &n
	}
	copy := NewNodeList(ns, n.ListMetadata.Continue)
	return &copy
}
