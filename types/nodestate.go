package types

import (
	"fmt"
	"strings"

	"github.com/slok/ragnarok/api/cluster/v1"
	pbns "github.com/slok/ragnarok/grpc/nodestatus"
)

// NodeStateParser has the required methods to transform node state.
type NodeStateParser interface {
	StrToNodeState(state string) (v1.NodeState, error)
	PBToNodeState(state pbns.State) (v1.NodeState, error)
	NodeStateToPB(state v1.NodeState) (pbns.State, error)
}

// nodeStateParser will convert node state in different types.
type nodeStateParser struct{}

// NodeStateTransformer is the utility to transform node status kinds.
var NodeStateTransformer = &nodeStateParser{}

// ParseNodeStateStr parses an string and returns a NodeState.
func (nodeStateParser) StrToNodeState(state string) (v1.NodeState, error) {
	switch strings.ToLower(state) {
	case "ready":
		return v1.ReadyNodeState, nil
	case "attacking":
		return v1.AttackingNodeState, nil
	case "reverting":
		return v1.RevertingNodeState, nil
	case "errored":
		return v1.ErroredNodeState, nil
	case "unknown":
		return v1.UnknownNodeState, nil
	default:
		return v1.UnknownNodeState, fmt.Errorf("invalid node state: %s", state)
	}
}

// ParseNodeStatePB parses a proto buffer node state and returns a NodeState.
func (nodeStateParser) PBToNodeState(state pbns.State) (v1.NodeState, error) {
	switch state {
	case pbns.State_READY:
		return v1.ReadyNodeState, nil
	case pbns.State_ATTACKING:
		return v1.AttackingNodeState, nil
	case pbns.State_REVERTING:
		return v1.RevertingNodeState, nil
	case pbns.State_ERRORED:
		return v1.ErroredNodeState, nil
	case pbns.State_UNKNOWN:
		return v1.UnknownNodeState, nil
	default:
		return v1.UnknownNodeState, fmt.Errorf("invalid node state: %s", state)
	}
}

// NodeStateToPB parses a node state to proto buffer and returns a PB node state.
func (nodeStateParser) NodeStateToPB(state v1.NodeState) (pbns.State, error) {
	switch state {
	case v1.ReadyNodeState:
		return pbns.State_READY, nil
	case v1.AttackingNodeState:
		return pbns.State_ATTACKING, nil
	case v1.RevertingNodeState:
		return pbns.State_REVERTING, nil
	case v1.ErroredNodeState:
		return pbns.State_ERRORED, nil
	case v1.UnknownNodeState:
		return pbns.State_UNKNOWN, nil
	default:
		return pbns.State_UNKNOWN, fmt.Errorf("invalid node state: %s", state)
	}
}
