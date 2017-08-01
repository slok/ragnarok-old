package client_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/failure"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	mfailure "github.com/slok/ragnarok/mocks/failure"
	mpbfs "github.com/slok/ragnarok/mocks/grpc/failurestatus"
	"github.com/slok/ragnarok/node/client"
	"github.com/slok/ragnarok/types"
)

func TestGetFailure(t *testing.T) {
	tests := []struct {
		name        string
		failure     *pbfs.Failure
		expFailure  *failure.Failure
		expRPCErr   bool // Expect GRPC call error.
		expTransErr bool // Expect transformation error.
	}{
		{
			name: "Get a failure correctly",
			failure: &pbfs.Failure{
				Id:            "test1",
				NodeID:        "node1",
				Definition:    "attacks:\n- attack1:\n    size: 524288000\n",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: pbfs.State_DISABLED,
			},
			expFailure: &failure.Failure{
				ID:     "test1",
				NodeID: "node1",
				Definition: failure.Definition{
					Attacks: []failure.AttackMap{
						{
							"attack1": attack.Opts{
								"size": 524288000,
							},
						},
					},
				},
				CurrentState:  types.EnabledFailureState,
				ExpectedState: types.DisabledFailureState,
			},
			expRPCErr:   false,
			expTransErr: false,
		},
		{
			name:        "RPC call failed",
			failure:     &pbfs.Failure{},
			expFailure:  &failure.Failure{},
			expRPCErr:   true,
			expTransErr: false,
		},
		{
			name:        "RPC call succesful but PB transformation error",
			failure:     &pbfs.Failure{},
			expFailure:  &failure.Failure{},
			expRPCErr:   false,
			expTransErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			var rpcErr, transErr error
			if test.expRPCErr {
				rpcErr = errors.New("wanted failure")
			}
			if test.expTransErr {
				transErr = errors.New("wanted failure")
			}

			// Create mocks.
			mc := &mpbfs.FailureStatusClient{}
			mc.On("GetFailure", mock.Anything, &pbfs.FailureId{Id: test.failure.GetId()}).Once().Return(test.failure, rpcErr)
			mp := &mfailure.Parser{}
			mp.On("PBToFailure", test.failure).Return(test.expFailure, transErr)

			// Create the service
			c, err := client.NewFailureGRPC(mc, mp, types.FailureStateTransformer, clock.Base(), log.Dummy)
			assert.NoError(err)

			// Make the call.
			f, err := c.GetFailure(test.failure.GetId())

			// Check.
			if test.expRPCErr || test.expTransErr {
				assert.Error(err)
			} else {
				if assert.NoError(err) {
					assert.Equal(test.expFailure, f)
				}
			}
			mc.AssertExpectations(t)
		})
	}
}
