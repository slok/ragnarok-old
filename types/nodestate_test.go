package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/types"
)

func TestNodeStateStringer(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		st    types.NodeState
		expSt string
	}{
		{types.ReadyNodeState, "ready"},
		{types.AttackingNodeState, "attacking"},
		{types.RevertingNodeState, "reverting"},
		{types.ErroredNodeState, "errored"},
		{types.UnknownNodeState, "unknown"},
		{99999, "unknown"},
	}

	for _, test := range tests {
		assert.Equal(test.expSt, test.st.String())
	}
}

func TestNodeStateParseStrToNS(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		st     string
		expSt  types.NodeState
		expErr bool
	}{
		{"ready", types.ReadyNodeState, false},
		{"attacking", types.AttackingNodeState, false},
		{"reverting", types.RevertingNodeState, false},
		{"errored", types.ErroredNodeState, false},
		{"unknown", types.UnknownNodeState, false},
		{"no-state", types.UnknownNodeState, true},
	}

	for _, test := range tests {
		gotSt, err := types.NodeStateTransformer.StrToNodeState(test.st)
		if test.expErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expSt, gotSt)
	}
}

func TestNodeStateParsePBToNS(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		st     pbns.State
		expSt  types.NodeState
		expErr bool
	}{
		{pbns.State_READY, types.ReadyNodeState, false},
		{pbns.State_ATTACKING, types.AttackingNodeState, false},
		{pbns.State_REVERTING, types.RevertingNodeState, false},
		{pbns.State_ERRORED, types.ErroredNodeState, false},
		{pbns.State_UNKNOWN, types.UnknownNodeState, false},
		{999999, types.UnknownNodeState, true},
	}

	for _, test := range tests {
		gotSt, err := types.NodeStateTransformer.PBToNodeState(test.st)
		if test.expErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expSt, gotSt)
	}
}

func TestNodeStateParseNSToPB(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		st     types.NodeState
		expSt  pbns.State
		expErr bool
	}{
		{types.ReadyNodeState, pbns.State_READY, false},
		{types.AttackingNodeState, pbns.State_ATTACKING, false},
		{types.RevertingNodeState, pbns.State_REVERTING, false},
		{types.ErroredNodeState, pbns.State_ERRORED, false},
		{types.UnknownNodeState, pbns.State_UNKNOWN, false},
		{999999, pbns.State_UNKNOWN, true},
	}

	for _, test := range tests {
		gotSt, err := types.NodeStateTransformer.NodeStateToPB(test.st)
		if test.expErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expSt, gotSt)
	}
}
