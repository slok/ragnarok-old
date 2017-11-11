package client_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	chaosv1pb "github.com/slok/ragnarok/api/chaos/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/clock"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	mpbfs "github.com/slok/ragnarok/mocks/grpc/failurestatus"
	mclient "github.com/slok/ragnarok/mocks/node/client"
	"github.com/slok/ragnarok/node/client"
	testpb "github.com/slok/ragnarok/test/pb"
)

func TestGetFailure(t *testing.T) {
	tests := []struct {
		name       string
		failure    *chaosv1.Failure
		expFailure *chaosv1.Failure
		expRPCErr  bool // Expect GRPC call error.
	}{
		{
			name: "Get a failure correctly",
			failure: &chaosv1.Failure{
				Metadata: chaosv1.FailureMetadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: chaosv1.FailureSpec{
					Attacks: []chaosv1.AttackMap{
						{
							"attack1": attack.Opts{
								"size": 524288000,
							},
						},
					},
				},
				Status: chaosv1.FailureStatus{
					CurrentState:  chaosv1.EnabledFailureState,
					ExpectedState: chaosv1.DisabledFailureState,
				},
			},
			expFailure: &chaosv1.Failure{
				TypeMeta: api.TypeMeta{
					Kind:    chaosv1.FailureKind,
					Version: chaosv1.FailureVersion,
				},
				Metadata: chaosv1.FailureMetadata{
					ID:     "test1",
					NodeID: "node1",
				},
				Spec: chaosv1.FailureSpec{
					Attacks: []chaosv1.AttackMap{
						{
							"attack1": attack.Opts{
								"size": float64(524288000),
							},
						},
					},
				},
				Status: chaosv1.FailureStatus{
					CurrentState:  chaosv1.EnabledFailureState,
					ExpectedState: chaosv1.DisabledFailureState,
				},
			},
			expRPCErr: false,
		},
		{
			name:       "RPC call failed",
			failure:    &chaosv1.Failure{},
			expFailure: &chaosv1.Failure{},
			expRPCErr:  true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			var rpcErr error
			if test.expRPCErr {
				rpcErr = errors.New("wanted failure")
			}

			// Create mocks.
			pbflr := testpb.CreatePBFailure(test.failure, t)
			mc := &mpbfs.FailureStatusClient{}
			mc.On("GetFailure", mock.Anything, mock.Anything).Once().Return(pbflr, rpcErr)

			// Create the service
			c, err := client.NewFailureGRPC(mc, serializer.PBSerializerDefault, clock.Base(), log.Dummy)
			require.NoError(err)

			// Make the call.
			f, err := c.GetFailure(test.failure.Metadata.ID)

			// Check.
			if test.expRPCErr {
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
		failures [][]*chaosv1.Failure
	}{
		{
			name:   "RPC call and stream correctly.",
			nodeID: "test1",
			failures: [][]*chaosv1.Failure{
				[]*chaosv1.Failure{
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id1"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.EnabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id2"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.EnabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id3"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.DisabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id4"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.DisabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id5"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.EnabledFailureState}},
				},
				[]*chaosv1.Failure{
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id6"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.DisabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id7"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.DisabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id8"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.EnabledFailureState}},
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
			for _, flrs := range test.failures {
				// Create the batches of failure pbs
				pbflrs := make([]*chaosv1pb.Failure, len(flrs))
				for i, flr := range flrs {
					pbflrs[i] = testpb.CreatePBFailure(flr, t)
				}
				fs := &pbfs.FailuresState{
					Failures: pbflrs,
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
			mc.On("FailureStateList", mock.Anything, mock.Anything).Once().Return(mstream, nil)

			// Create the service
			c, err := client.NewFailureGRPC(mc, serializer.PBSerializerDefault, clock.Base(), log.Dummy)
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
		failures [][]*chaosv1.Failure
	}{
		{
			name:   "RPC call and stream correctly.",
			nodeID: "test1",
			failures: [][]*chaosv1.Failure{
				[]*chaosv1.Failure{
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id1"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.EnabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id2"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.EnabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id3"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.DisabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id4"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.DisabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id5"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.EnabledFailureState}},
				},
				[]*chaosv1.Failure{
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id6"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.DisabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id7"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.DisabledFailureState}},
					&chaosv1.Failure{Metadata: chaosv1.FailureMetadata{ID: "id8"}, Status: chaosv1.FailureStatus{ExpectedState: chaosv1.EnabledFailureState}},
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
			for _, flrs := range test.failures {
				// Create the batches of failure pbs
				pbflrs := make([]*chaosv1pb.Failure, len(flrs))
				for i, flr := range flrs {
					pbflrs[i] = testpb.CreatePBFailure(flr, t)
				}
				fs := &pbfs.FailuresState{
					Failures: pbflrs,
				}
				mstream.On("Recv").Once().Return(fs, nil)
				// This will check correct handling of grpc results from the server.
				mfsh.On("ProcessFailureStates", flrs).Once().Return(nil)
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
			mc.On("FailureStateList", mock.Anything, mock.Anything).Return(mstream, nil)

			// Create the service
			stopC := make(chan struct{})
			c, err := client.NewFailureGRPC(mc, serializer.PBSerializerDefault, clock.Base(), log.Dummy)
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
