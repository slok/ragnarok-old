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
	ID     string // ID is the id of the node
	Master bool   // Master will telle what kind of node is.
}

// NodeSpec has the node specific fields.
type NodeSpec struct {
	Labels NodeLabels // Labels are the tags related with the node.
}

// NodeStatus has the state fo the node.
type NodeStatus struct {
	State    NodeState // State is the state of the Node
	Creation time.Time // Creation is when the creation of the node happenned.
}

// Node is an internal and simplified representation of a failure node on the masters
type Node struct {
	Metadata NodeMetadata
	Spec     NodeSpec
	Status   NodeStatus
}

// GetObjectKind satisfies Object interface.
func (n *Node) GetObjectKind() api.Kind {
	return NodeKind
}

// GetObjectVersion satisfies Object interface.
func (n *Node) GetObjectVersion() api.Version {
	return NodeVersion
}
