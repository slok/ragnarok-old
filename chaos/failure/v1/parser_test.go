package v1_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/chaos/failure/v1"
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
				Metadata: v1.Metadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: v1.Spec{
					Attacks: []v1.AttackMap{
						{
							"attack1": attack.Opts{
								"size": 524288000,
							},
						},
					},
				},
				Status: v1.Status{
					CurrentState:  types.EnabledFailureState,
					ExpectedState: types.DisabledFailureState,
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
			gotFailure, err := v1.Transformer.PBToFailure(test.failure)

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
				Metadata: v1.Metadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: v1.Spec{
					Attacks: []v1.AttackMap{
						{
							"attack1": attack.Opts{
								"size": 524288000,
							},
						},
					},
				},
				Status: v1.Status{
					CurrentState:  types.EnabledFailureState,
					ExpectedState: types.DisabledFailureState,
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
				Metadata: v1.Metadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: v1.Spec{},
				Status: v1.Status{
					CurrentState:  9999999999,
					ExpectedState: types.DisabledFailureState,
				},
			},
			expErr: true,
		},
		{
			name: "Error on expected state transformation",
			failure: &v1.Failure{
				Metadata: v1.Metadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: v1.Spec{},
				Status: v1.Status{

					CurrentState:  types.DisabledFailureState,
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
			gotFailure, err := v1.Transformer.FailureToPB(test.failure)

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
