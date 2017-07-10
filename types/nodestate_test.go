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

func TestNodeStateParseStr(t *testing.T) {
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
		gotSt, err := types.ParseNodeStateStr(test.st)
		if test.expErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expSt, gotSt)
	}
}

func TestNodeStateParsePB(t *testing.T) {
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
		gotSt, err := types.ParseNodeStatePB(test.st)
		if test.expErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expSt, gotSt)
	}
}
