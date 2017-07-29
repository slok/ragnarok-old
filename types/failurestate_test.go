package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/types"
)

func TestFailureStateStringer(t *testing.T) {
	tests := []struct {
		st    types.FailureState
		expSt string
	}{
		{types.EnabledFailureState, "enabled"},
		{types.ExecutingFailureState, "executing"},
		{types.RevertingFailureState, "reverting"},
		{types.DisabledFailureState, "disabled"},
		{types.ErroredFailureState, "errored"},
		{types.ErroredRevertingFailureState, "erroredreverting"},
		{types.UnknownFailureState, "unknown"},
		{99999, "unknown"},
	}

	for _, test := range tests {
		t.Run(test.expSt, func(t *testing.T) {
			assert := assert.New(t)
			assert.Equal(test.expSt, test.st.String())
		})
	}
}

func TestFailureStateParseStrToFS(t *testing.T) {
	tests := []struct {
		st     string
		expSt  types.FailureState
		expErr bool
	}{
		{"enabled", types.EnabledFailureState, false},
		{"executing", types.ExecutingFailureState, false},
		{"reverting", types.RevertingFailureState, false},
		{"disabled", types.DisabledFailureState, false},
		{"errored", types.ErroredFailureState, false},
		{"erroredreverting", types.ErroredRevertingFailureState, false},
		{"no-state", types.UnknownFailureState, true},
	}

	for _, test := range tests {
		t.Run(test.st, func(t *testing.T) {
			assert := assert.New(t)

			gotSt, err := types.FailureStateTransformer.StrToFailureState(test.st)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			assert.Equal(test.expSt, gotSt)
		})
	}
}

func TestFailureStateParsePBToFS(t *testing.T) {
	tests := []struct {
		st     pbfs.State
		expSt  types.FailureState
		expErr bool
	}{
		{pbfs.State_ENABLED, types.EnabledFailureState, false},
		{pbfs.State_EXECUTING, types.ExecutingFailureState, false},
		{pbfs.State_REVERTING, types.RevertingFailureState, false},
		{pbfs.State_DISABLED, types.DisabledFailureState, false},
		{pbfs.State_ERRORED, types.ErroredFailureState, false},
		{pbfs.State_ERRORED_REVERTING, types.ErroredRevertingFailureState, false},
		{pbfs.State_UNKNOWN, types.UnknownFailureState, true},
		{999999, types.UnknownFailureState, true},
	}

	for _, test := range tests {
		t.Run(test.st.String(), func(t *testing.T) {
			assert := assert.New(t)

			gotSt, err := types.FailureStateTransformer.PBToFailureState(test.st)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			assert.Equal(test.expSt, gotSt)
		})
	}
}

func TestFailureStateParseFSToPB(t *testing.T) {
	tests := []struct {
		st     types.FailureState
		expSt  pbfs.State
		expErr bool
	}{
		{types.EnabledFailureState, pbfs.State_ENABLED, false},
		{types.ExecutingFailureState, pbfs.State_EXECUTING, false},
		{types.RevertingFailureState, pbfs.State_REVERTING, false},
		{types.DisabledFailureState, pbfs.State_DISABLED, false},
		{types.ErroredFailureState, pbfs.State_ERRORED, false},
		{types.ErroredRevertingFailureState, pbfs.State_ERRORED_REVERTING, false},
		{types.UnknownFailureState, pbfs.State_UNKNOWN, true},
		{999999, pbfs.State_UNKNOWN, true},
	}

	for _, test := range tests {
		t.Run(test.st.String(), func(t *testing.T) {
			assert := assert.New(t)

			gotSt, err := types.FailureStateTransformer.FailureStateToPB(test.st)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			assert.Equal(test.expSt, gotSt)
		})
	}
}
