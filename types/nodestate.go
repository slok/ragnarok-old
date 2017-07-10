package types

import (
	"fmt"
	"strings"
)

// NodeState is the reprensetation of the node state
type NodeState int

const (
	// UnknownNodeState is an unknown status
	UnknownNodeState NodeState = iota
	// ReadyNodeState is the state when a node is ready for accepting attacks
	ReadyNodeState
	// AttackingNodeState is the state when a node is aplying an attack
	AttackingNodeState
	// RevertingNodeState is the state when a node is reverting an applied attack
	RevertingNodeState
	// ErroredNodeState is the state when a node is in error state
	ErroredNodeState
)

// String implements the stringer interface
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

// ParseNodeState parses an string and returns a NodeState
func ParseNodeState(state string) (NodeState, error) {
	switch strings.ToLower(state) {
	case "ready":
		return ReadyNodeState, nil
	case "attacking":
		return AttackingNodeState, nil
	case "reverting":
		return RevertingNodeState, nil
	case "errored":
		return ErroredNodeState, nil
	case "unknown":
		return UnknownNodeState, nil
	default:
		return UnknownNodeState, fmt.Errorf("invalid node state: %s", state)
	}
}
