package types_test

import (
	"testing"

	"github.com/slok/ragnarok/types"
	"github.com/stretchr/testify/assert"
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

func TestNodeStateParse(t *testing.T) {
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
		gotSt, err := types.ParseNodeState(test.st)
		if test.expErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expSt, gotSt)
	}
}
