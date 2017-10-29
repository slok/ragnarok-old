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
)

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

// NodeLabels is a key value pair map.
type NodeLabels map[string]string

// NodeMetadata has the node metadata fields
type NodeMetadata struct {
	ID     string `json:"id,omitempty"`     // ID is the id of the node
	Master bool   `json:"master,omitempty"` // Master will telle what kind of node is.
}

// NodeSpec has the node specific fields.
type NodeSpec struct {
	Labels NodeLabels `json:"labels,omitempty"` // Labels are the tags related with the node.
}

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

	Metadata NodeMetadata `json:"metadata,omitempty"`
	Spec     NodeSpec     `json:"spec,omitempty"`
	Status   NodeStatus   `json:"status,omitempty"`
}
