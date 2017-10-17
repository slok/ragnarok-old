package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api/cluster/v1"
	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/types"
)

func TestNodeStateStringer(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		st    v1.NodeState
		expSt string
	}{
		{v1.ReadyNodeState, "ready"},
		{v1.AttackingNodeState, "attacking"},
		{v1.RevertingNodeState, "reverting"},
		{v1.ErroredNodeState, "errored"},
		{v1.UnknownNodeState, "unknown"},
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
		expSt  v1.NodeState
		expErr bool
	}{
		{"ready", v1.ReadyNodeState, false},
		{"attacking", v1.AttackingNodeState, false},
		{"reverting", v1.RevertingNodeState, false},
		{"errored", v1.ErroredNodeState, false},
		{"unknown", v1.UnknownNodeState, false},
		{"no-state", v1.UnknownNodeState, true},
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
		expSt  v1.NodeState
		expErr bool
	}{
		{pbns.State_READY, v1.ReadyNodeState, false},
		{pbns.State_ATTACKING, v1.AttackingNodeState, false},
		{pbns.State_REVERTING, v1.RevertingNodeState, false},
		{pbns.State_ERRORED, v1.ErroredNodeState, false},
		{pbns.State_UNKNOWN, v1.UnknownNodeState, false},
		{999999, v1.UnknownNodeState, true},
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
		st     v1.NodeState
		expSt  pbns.State
		expErr bool
	}{
		{v1.ReadyNodeState, pbns.State_READY, false},
		{v1.AttackingNodeState, pbns.State_ATTACKING, false},
		{v1.RevertingNodeState, pbns.State_REVERTING, false},
		{v1.ErroredNodeState, pbns.State_ERRORED, false},
		{v1.UnknownNodeState, pbns.State_UNKNOWN, false},
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
