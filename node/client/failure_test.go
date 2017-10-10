package client_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/chaos/failure"
	"github.com/slok/ragnarok/clock"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	mfailure "github.com/slok/ragnarok/mocks/chaos/failure"
	mpbfs "github.com/slok/ragnarok/mocks/grpc/failurestatus"
	mclient "github.com/slok/ragnarok/mocks/node/client"
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
			name:        "RPC call succesful but PB result transformation error",
			failure:     &pbfs.Failure{},
			expFailure:  &failure.Failure{},
			expRPCErr:   false,
			expTransErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

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
			require.NoError(err)

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

func TestFailureStateListStreamingOK(t *testing.T) {
	tests := []struct {
		name     string
		nodeID   string
		statuses [][]*pbfs.Failure
	}{
		{
			name:   "RPC call and stream correctly.",
			nodeID: "test1",
			statuses: [][]*pbfs.Failure{
				[]*pbfs.Failure{
					&pbfs.Failure{Id: "id1", ExpectedState: pbfs.State_ENABLED},
					&pbfs.Failure{Id: "id2", ExpectedState: pbfs.State_ENABLED},
					&pbfs.Failure{Id: "id3", ExpectedState: pbfs.State_DISABLED},
					&pbfs.Failure{Id: "id4", ExpectedState: pbfs.State_DISABLED},
					&pbfs.Failure{Id: "id5", ExpectedState: pbfs.State_ENABLED},
				},
				[]*pbfs.Failure{
					&pbfs.Failure{Id: "id6", ExpectedState: pbfs.State_DISABLED},
					&pbfs.Failure{Id: "id7", ExpectedState: pbfs.State_DISABLED},
					&pbfs.Failure{Id: "id8", ExpectedState: pbfs.State_ENABLED},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// Used to track the goroutine call finished.
			finishedC := make(chan struct{})

			// Create mocks.
			// Mock the streaming of batches from the server.
			// Mock the  handler of the states in order to verify the proper handling of the stream statuses.
			mstream := &mpbfs.FailureStatus_FailureStateListClient{}
			mfsh := &mclient.FailureStateHandler{}
			mstream.On("Context").Return(context.Background())
			for _, st := range test.statuses {
				fs := &pbfs.FailuresState{
					Failures: st,
				}
				mstream.On("Recv").Once().Return(fs, nil)
				mfsh.On("ProcessFailureStates", mock.Anything).Once().Return(nil)
			}
			// Ignore next streaming read receive calls.
			mstream.On("Recv").Return(&pbfs.FailuresState{}, nil).Run(func(args mock.Arguments) {
				finishedC <- struct{}{}
			})

			// Mock the server GRPC real call.
			mc := &mpbfs.FailureStatusClient{}
			mc.On("FailureStateList", mock.Anything, &pbfs.NodeId{Id: test.nodeID}).Once().Return(mstream, nil)

			// Create the service
			c, err := client.NewFailureGRPC(mc, failure.Transformer, types.FailureStateTransformer, clock.Base(), log.Dummy)
			require.NoError(err)
			err = c.ProcessFailureStateStreaming(test.nodeID, mfsh, nil)
			if assert.NoError(err) {
				// Wait to the stream activity.
				select {
				case <-time.After(10 * time.Millisecond):
					assert.Fail("timeout waiting to receive data from the server stream")
				case <-finishedC:
				}

				// Check the correct handling of the states.
				mfsh.AssertExpectations(t)
			}
		})
	}
}

func TestFailureStateListStreamingOKWithStop(t *testing.T) {
	tests := []struct {
		name     string
		nodeID   string
		statuses [][]*pbfs.Failure
	}{
		{
			name:   "RPC call and stream correctly.",
			nodeID: "test1",
			statuses: [][]*pbfs.Failure{
				[]*pbfs.Failure{
					&pbfs.Failure{Id: "id1", ExpectedState: pbfs.State_ENABLED},
					&pbfs.Failure{Id: "id2", ExpectedState: pbfs.State_ENABLED},
					&pbfs.Failure{Id: "id3", ExpectedState: pbfs.State_DISABLED},
					&pbfs.Failure{Id: "id4", ExpectedState: pbfs.State_DISABLED},
					&pbfs.Failure{Id: "id5", ExpectedState: pbfs.State_ENABLED},
				},
				[]*pbfs.Failure{
					&pbfs.Failure{Id: "id6", ExpectedState: pbfs.State_DISABLED},
					&pbfs.Failure{Id: "id7", ExpectedState: pbfs.State_DISABLED},
					&pbfs.Failure{Id: "id8", ExpectedState: pbfs.State_ENABLED},
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// Used to track the goroutine call finished.
			finishedC := make(chan struct{})

			// Create mocks.
			// Mock the streaming of batches from the server.
			// Mock the  handler of the states in order to verify the proper handling of the stream statuses.
			mstream := &mpbfs.FailureStatus_FailureStateListClient{}
			mfsh := &mclient.FailureStateHandler{}
			mstream.On("Context").Return(context.Background())
			for _, st := range test.statuses {
				fs := &pbfs.FailuresState{
					Failures: st,
				}
				mstream.On("Recv").Once().Return(fs, nil)
				mfsh.On("ProcessFailureStates", mock.Anything).Once().Return(nil)
			}
			// Ignore next streaming read receive calls.
			mstream.On("Recv").Return(&pbfs.FailuresState{}, nil).Run(func(args mock.Arguments) {
				select {
				case finishedC <- struct{}{}:
				default:
				}
			})

			// Mock the server GRPC real call.
			mc := &mpbfs.FailureStatusClient{}
			mc.On("FailureStateList", mock.Anything, &pbfs.NodeId{Id: test.nodeID}).Return(mstream, nil)

			// Create the service
			stopC := make(chan struct{})
			c, err := client.NewFailureGRPC(mc, failure.Transformer, types.FailureStateTransformer, clock.Base(), log.Dummy)
			require.NoError(err)
			err = c.ProcessFailureStateStreaming(test.nodeID, mfsh, stopC)
			if assert.NoError(err) {
				// Wait to the stream activity.
				select {
				case <-time.After(10 * time.Millisecond):
					assert.Fail("timeout waiting to receive data from the server stream")
				case <-finishedC:
				}

				// Check the correct handling of the states.
				mfsh.AssertExpectations(t)

				// Check stop is ok.
				err := c.ProcessFailureStateStreaming(test.nodeID, mfsh, stopC)
				require.Error(err)
				stopC <- struct{}{}
				time.Sleep(5 * time.Millisecond)
				err = c.ProcessFailureStateStreaming(test.nodeID, mfsh, stopC)
				assert.NoError(err)
			}
		})
	}
}
