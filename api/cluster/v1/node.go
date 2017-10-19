package v1

import (
	"github.com/slok/ragnarok/api"
)

const (
	// NodeKind is the kind a failure.
	NodeKind = "cluster/v1/node"
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

// Node is an internal and simplified representation of a failure node on the masters
// TODO: Rethink the reuse of node
type Node struct {
	ID     string     // ID is the id of the node
	Labels NodeLabels // Labels are the tags related with the node
	State  NodeState  // State is the state of the Node
}

// GetObjectKind satisfies Object interface.
func (n *Node) GetObjectKind() api.Kind {
	return NodeKind
}
