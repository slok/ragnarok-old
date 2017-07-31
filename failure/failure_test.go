package failure_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/failure"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/types"
)

func TestFailureTrasformPB2Failure(t *testing.T) {

	tests := []struct {
		name       string
		failure    *pbfs.Failure
		expFailure *failure.Failure
		expErr     bool
	}{
		{
			name: "Simple conversion",
			failure: &pbfs.Failure{
				Id:            "test1",
				NodeID:        "node1",
				Definition:    "{}\n",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: pbfs.State_DISABLED,
			},
			expFailure: &failure.Failure{
				ID:            "test1",
				NodeID:        "node1",
				Definition:    failure.Definition{},
				CurrentState:  types.EnabledFailureState,
				ExpectedState: types.DisabledFailureState,
			},
		},
		{
			name: "Error on definition unmarshaling",
			failure: &pbfs.Failure{
				Id:            "test1",
				NodeID:        "node1",
				Definition:    "{{{",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: pbfs.State_DISABLED,
			},
			expErr: true,
		},
		{
			name: "Error on current state transformation",
			failure: &pbfs.Failure{
				Id:            "test1",
				NodeID:        "node1",
				Definition:    "{}\n",
				CurrentState:  9999,
				ExpectedState: pbfs.State_DISABLED,
			},
			expErr: true,
		},
		{
			name: "Error on expected state transformation",
			failure: &pbfs.Failure{
				Id:            "test1",
				NodeID:        "node1",
				Definition:    "{}\n",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: 9999,
			},
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Transform.
			gotFailure, err := failure.Transformer.PBToFailure(test.failure)

			// Check.
			if test.expErr {
				assert.Error(err)
			} else {
				if assert.NoError(err) {
					assert.Equal(test.expFailure, gotFailure)
				}
			}
		})
	}
}
