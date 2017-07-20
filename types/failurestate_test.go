package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
