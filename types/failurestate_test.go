package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/types"
)

func TestFailureStateStringer(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		st    types.FailureState
		expSt string
	}{
		{types.EnabledFailureState, "enabled"},
		{types.RevertingFailureState, "reverting"},
		{types.DisabledFailureState, "disabled"},
		{types.UnknownFailureState, "unknown"},
		{99999, "unknown"},
	}

	for _, test := range tests {
		assert.Equal(test.expSt, test.st.String())
	}
}

func TestFailureStateParseStrToFS(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		st     string
		expSt  types.FailureState
		expErr bool
	}{
		{"enabled", types.EnabledFailureState, false},
		{"reverting", types.RevertingFailureState, false},
		{"disabled", types.DisabledFailureState, false},
		{"no-state", types.UnknownFailureState, true},
	}

	for _, test := range tests {
		gotSt, err := types.FailureStateTransformer.StrToFailureState(test.st)
		if test.expErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expSt, gotSt)
	}
}

func TestFailureStateParsePBToFS(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		st     pbfs.State
		expSt  types.FailureState
		expErr bool
	}{
		{pbfs.State_ENABLED, types.EnabledFailureState, false},
		{pbfs.State_REVERTING, types.RevertingFailureState, false},
		{pbfs.State_DISABLED, types.DisabledFailureState, false},
		{pbfs.State_UNKNOWN, types.UnknownFailureState, true},
		{999999, types.UnknownFailureState, true},
	}

	for _, test := range tests {
		gotSt, err := types.FailureStateTransformer.PBToFailureState(test.st)
		if test.expErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expSt, gotSt)
	}
}

func TestFailureStateParseFSToPB(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		st     types.FailureState
		expSt  pbfs.State
		expErr bool
	}{
		{types.EnabledFailureState, pbfs.State_ENABLED, false},
		{types.RevertingFailureState, pbfs.State_REVERTING, false},
		{types.DisabledFailureState, pbfs.State_DISABLED, false},
		{types.UnknownFailureState, pbfs.State_UNKNOWN, true},
		{999999, pbfs.State_UNKNOWN, true},
	}

	for _, test := range tests {
		gotSt, err := types.FailureStateTransformer.FailureStateToPB(test.st)
		if test.expErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		assert.Equal(test.expSt, gotSt)
	}
}
