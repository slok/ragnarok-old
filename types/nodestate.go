package types

import (
	"fmt"
	"strings"

	pbns "github.com/slok/ragnarok/grpc/nodestatus"
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

// NodeStateParser has the required methods to transform node state.
type NodeStateParser interface {
	StrToNodeState(state string) (NodeState, error)
	PBToNodeState(state pbns.State) (NodeState, error)
	NodeStateToPB(state NodeState) (pbns.State, error)
}

// nodeStateParser will convert node state in different types.
type nodeStateParser struct{}

// NodeStateTransformer is the utility to transform node status kinds.
var NodeStateTransformer = &nodeStateParser{}

// ParseNodeStateStr parses an string and returns a NodeState.
func (nodeStateParser) StrToNodeState(state string) (NodeState, error) {
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

// ParseNodeStatePB parses a proto buffer node state and returns a NodeState.
func (nodeStateParser) PBToNodeState(state pbns.State) (NodeState, error) {
	switch state {
	case pbns.State_READY:
		return ReadyNodeState, nil
	case pbns.State_ATTACKING:
		return AttackingNodeState, nil
	case pbns.State_REVERTING:
		return RevertingNodeState, nil
	case pbns.State_ERRORED:
		return ErroredNodeState, nil
	case pbns.State_UNKNOWN:
		return UnknownNodeState, nil
	default:
		return UnknownNodeState, fmt.Errorf("invalid node state: %s", state)
	}
}

// NodeStateToPB parses a node state to proto buffer and returns a PB node state.
func (nodeStateParser) NodeStateToPB(state NodeState) (pbns.State, error) {
	switch state {
	case ReadyNodeState:
		return pbns.State_READY, nil
	case AttackingNodeState:
		return pbns.State_ATTACKING, nil
	case RevertingNodeState:
		return pbns.State_REVERTING, nil
	case ErroredNodeState:
		return pbns.State_ERRORED, nil
	case UnknownNodeState:
		return pbns.State_UNKNOWN, nil
	default:
		return pbns.State_UNKNOWN, fmt.Errorf("invalid node state: %s", state)
	}
}
