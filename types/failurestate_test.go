package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api/chaos/v1"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/types"
)

func TestFailureStateParseStrToFS(t *testing.T) {
	tests := []struct {
		st     string
		expSt  v1.FailureState
		expErr bool
	}{
		{"enabled", v1.EnabledFailureState, false},
		{"executing", v1.ExecutingFailureState, false},
		{"reverting", v1.RevertingFailureState, false},
		{"disabled", v1.DisabledFailureState, false},
		{"stale", v1.StaleFailureState, false},
		{"errored", v1.ErroredFailureState, false},
		{"erroredreverting", v1.ErroredRevertingFailureState, false},
		{"no-state", v1.UnknownFailureState, true},
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
		expSt  v1.FailureState
		expErr bool
	}{
		{pbfs.State_ENABLED, v1.EnabledFailureState, false},
		{pbfs.State_EXECUTING, v1.ExecutingFailureState, false},
		{pbfs.State_REVERTING, v1.RevertingFailureState, false},
		{pbfs.State_DISABLED, v1.DisabledFailureState, false},
		{pbfs.State_STALE, v1.StaleFailureState, false},
		{pbfs.State_ERRORED, v1.ErroredFailureState, false},
		{pbfs.State_ERRORED_REVERTING, v1.ErroredRevertingFailureState, false},
		{pbfs.State_UNKNOWN, v1.UnknownFailureState, true},
		{999999, v1.UnknownFailureState, true},
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
		st     v1.FailureState
		expSt  pbfs.State
		expErr bool
	}{
		{v1.EnabledFailureState, pbfs.State_ENABLED, false},
		{v1.ExecutingFailureState, pbfs.State_EXECUTING, false},
		{v1.RevertingFailureState, pbfs.State_REVERTING, false},
		{v1.DisabledFailureState, pbfs.State_DISABLED, false},
		{v1.StaleFailureState, pbfs.State_STALE, false},
		{v1.ErroredFailureState, pbfs.State_ERRORED, false},
		{v1.ErroredRevertingFailureState, pbfs.State_ERRORED_REVERTING, false},
		{v1.UnknownFailureState, pbfs.State_UNKNOWN, true},
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
