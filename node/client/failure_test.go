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
			c, err := client.NewFailureGRPC(0, mc, mp, types.FailureStateTransformer, clock.Base(), log.Dummy)
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
		name        string
		nodeID      string
		statuses    []map[string][]string
		expStatuses map[string][]string
	}{
		{
			name:   "RPC call and stream correctly.",
			nodeID: "test1",
			statuses: []map[string][]string{
				map[string][]string{
					"enabled":  []string{"id1", "id2", "id5"},
					"disabled": []string{"id3", "id4"},
				},
				map[string][]string{
					"enabled":  []string{"id6", "id8"},
					"disabled": []string{"id7"},
				},
			},
			expStatuses: map[string][]string{
				"enabled":  []string{"id1", "id2", "id5", "id6", "id8"},
				"disabled": []string{"id3", "id4", "id7"},
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
			mstream := &mpbfs.FailureStatus_FailureStateListClient{}
			mstream.On("Context").Return(context.Background())
			for _, st := range test.statuses {
				fs := &pbfs.FailuresExpectedState{
					EnabledFailureId:  st["enabled"],
					DisabledFailureId: st["disabled"],
				}
				mstream.On("Recv").Once().Return(fs, nil)
			}
			mstream.On("Recv").Return(&pbfs.FailuresExpectedState{}, nil).Run(func(args mock.Arguments) {
				finishedC <- struct{}{}
			}) // Ignore next receive calls.
			mc := &mpbfs.FailureStatusClient{}
			mc.On("FailureStateList", mock.Anything, &pbfs.NodeId{Id: test.nodeID}).Once().Return(mstream, nil)
			mp := &mfailure.Parser{}

			// Create the service
			c, err := client.NewFailureGRPC(5, mc, mp, types.FailureStateTransformer, clock.Base(), log.Dummy)
			require.NoError(err)
			esC, dsC, err := c.FailureStateList(test.nodeID)
			if assert.NoError(err) {

				// Wait to the stream activity.
				select {
				case <-time.After(10 * time.Millisecond):
					assert.Fail("timeout waiting to receive data from the server stream")
				case <-finishedC:
				}

				// Get the data from our channels and check.
				gotEnabled := []string{}
				gotDisabled := []string{}

				for {
					select {
					case st := <-esC:
						gotEnabled = append(gotEnabled, st)
					case st := <-dsC:
						gotDisabled = append(gotDisabled, st)
					case <-time.After(10 * time.Millisecond): // Give some time to process all the messages.
						if assert.NoError(err) {
							gotStatuses := map[string][]string{
								"enabled":  gotEnabled,
								"disabled": gotDisabled,
							}
							assert.Equal(test.expStatuses, gotStatuses)
						}
						return
					}
				}
			}

		})
	}
}

func TestFailureStateListStreamingContextDone(t *testing.T) {
	tests := []struct {
		name     string
		nodeID   string
		errClose bool
		expErr   bool
	}{
		{
			name:     "GRPC stream closed the context, client closes connection ok",
			nodeID:   "testnode1",
			errClose: false,
			expErr:   false,
		},
		{
			name:     "GRPC stream closed the context, error when client closes connection",
			nodeID:   "testnode1",
			errClose: true,
			expErr:   false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			var errClose error
			if test.errClose {
				errClose = errors.New("wanted error")
			}

			finishedTestC := make(chan struct{}) // Used to know when the tests has finished.

			// Create mocks.
			mstream := &mpbfs.FailureStatus_FailureStateListClient{}
			ctx, ccl := context.WithCancel(context.Background())
			ccl()
			mstream.On("Context").Once().Return(ctx)
			mstream.On("CloseSend").Once().Return(errClose).Run(func(args mock.Arguments) {
				// Send the signal of tests is finished.
				finishedTestC <- struct{}{}
			})
			mc := &mpbfs.FailureStatusClient{}
			mc.On("FailureStateList", mock.Anything, &pbfs.NodeId{Id: test.nodeID}).Once().Return(mstream, nil)
			mp := &mfailure.Parser{}

			// Create the service
			c, err := client.NewFailureGRPC(0, mc, mp, types.FailureStateTransformer, clock.Base(), log.Dummy)
			require.NoError(err)
			_, _, err = c.FailureStateList(test.nodeID)
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			// Wait to finish the tests.
			select {
			case <-time.After(10 * time.Millisecond):
				assert.Fail("timeout waiting to finish tets")
			case <-finishedTestC:
			}
			mstream.AssertExpectations(t)
		})
	}
}
