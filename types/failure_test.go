package types_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/attack"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/types"
)

func TestFailureTrasformPB2Failure(t *testing.T) {

	tests := []struct {
		name       string
		failure    *pbfs.Failure
		expFailure *v1.Failure
		expErr     bool
	}{
		{
			name: "Simple conversion",
			failure: &pbfs.Failure{
				Id:            "test1",
				NodeID:        "node1",
				Definition:    "attacks:\n- attack1:\n    size: 524288000\n",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: pbfs.State_DISABLED,
			},
			expFailure: &v1.Failure{
				Metadata: v1.FailureMetadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: v1.FailureSpec{
					Attacks: []v1.AttackMap{
						{
							"attack1": attack.Opts{
								"size": 524288000,
							},
						},
					},
				},
				Status: v1.FailureStatus{
					CurrentState:  v1.EnabledFailureState,
					ExpectedState: v1.DisabledFailureState,
				},
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
			gotFailure, err := types.FailureTransformer.PBToFailure(test.failure)

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

func TestFailureTrasformFailure2PB(t *testing.T) {

	tests := []struct {
		name       string
		failure    *v1.Failure
		expFailure *pbfs.Failure
		expErr     bool
	}{
		{
			name: "Simple conversion",
			failure: &v1.Failure{
				Metadata: v1.FailureMetadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: v1.FailureSpec{
					Attacks: []v1.AttackMap{
						{
							"attack1": attack.Opts{
								"size": 524288000,
							},
						},
					},
				},
				Status: v1.FailureStatus{
					CurrentState:  v1.EnabledFailureState,
					ExpectedState: v1.DisabledFailureState,
				},
			},
			expFailure: &pbfs.Failure{
				Id:            "test1",
				NodeID:        "node1",
				Definition:    "attacks:\n- attack1:\n    size: 524288000\n",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: pbfs.State_DISABLED,
			},
		},
		{
			name: "Error on current state transformation",
			failure: &v1.Failure{
				Metadata: v1.FailureMetadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: v1.FailureSpec{},
				Status: v1.FailureStatus{
					CurrentState:  9999999999,
					ExpectedState: v1.DisabledFailureState,
				},
			},
			expErr: true,
		},
		{
			name: "Error on expected state transformation",
			failure: &v1.Failure{
				Metadata: v1.FailureMetadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: v1.FailureSpec{},
				Status: v1.FailureStatus{

					CurrentState:  v1.DisabledFailureState,
					ExpectedState: 9999999999,
				},
			},
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Transform.
			gotFailure, err := types.FailureTransformer.FailureToPB(test.failure)

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
