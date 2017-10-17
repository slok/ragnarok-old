package server_test

import (
	"errors"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context" // TODO: Change when GRPC supports std librarie context

	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/clock"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/server"
	mclock "github.com/slok/ragnarok/mocks/clock"
	mservice "github.com/slok/ragnarok/mocks/service"
	tgrpc "github.com/slok/ragnarok/test/grpc"
	"github.com/slok/ragnarok/types"
)

func TestMasterGRPCServiceServerRegisterNode(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	tests := []struct {
		id        string
		tags      map[string]string
		shouldErr bool
	}{
		{"test1", nil, false},
		{"test1", map[string]string{"address": "10.234.012"}, true},
		{"test1", map[string]string{"address": "10.234.013", "kind": "complex"}, false},
	}

	for _, test := range tests {
		// Create service mocks.
		mfss := &mservice.FailureStatusService{}
		mnss := &mservice.NodeStatusService{}
		var expErr error
		if test.shouldErr {
			expErr = errors.New("wanted error")
		}
		mnss.On("Register", test.id, test.tags).Once().Return(expErr)

		// Create our server.
		l, err := net.Listen("tcp", "127.0.0.1:0") // :0 for a random port.
		require.NoError(err)
		defer l.Close()
		s := server.NewMasterGRPCServiceServer(mfss, mnss, l, clock.Base(), log.Dummy)
		// Serve in background.
		go func() {
			s.Serve()
		}()

		// Create our client.
		testCli, err := tgrpc.NewTestClient(l.Addr().String())
		require.NoError(err)
		defer testCli.Close()

		// Make call.
		n := &pbns.Node{
			Id:   test.id,
			Tags: test.tags,
		}
		_, err = testCli.NodeStatusRegister(context.Background(), n)

		// Check.
		if test.shouldErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		// Assert correct calls on our logic.
		mnss.AssertExpectations(t)
	}
}

func TestMasterGRPCServiceServerNodeHeartbeat(t *testing.T) {
	require := require.New(t)
	assert := assert.New(t)

	tests := []struct {
		id        string
		state     pbns.State
		expState  clusterv1.NodeState
		shouldErr bool
	}{
		{"test1", pbns.State_READY, clusterv1.ReadyNodeState, false},
		{"test1", pbns.State_UNKNOWN, clusterv1.UnknownNodeState, false},
		{"test1", pbns.State_ERRORED, clusterv1.ErroredNodeState, false},
		{"test1", pbns.State_ATTACKING, clusterv1.AttackingNodeState, false},
		{"test1", pbns.State_REVERTING, clusterv1.RevertingNodeState, false},
		{"test1", pbns.State_REVERTING, clusterv1.RevertingNodeState, true},
	}

	for _, test := range tests {
		// Create service mocks.
		mfss := &mservice.FailureStatusService{}
		mnss := &mservice.NodeStatusService{}
		var expErr error
		if test.shouldErr {
			expErr = errors.New("wanted error")
		}
		mnss.On("Heartbeat", test.id, test.expState).Once().Return(expErr)

		// Create our server.
		l, err := net.Listen("tcp", "127.0.0.1:0") // :0 for a random port.
		require.NoError(err)
		defer l.Close()
		s := server.NewMasterGRPCServiceServer(mfss, mnss, l, clock.Base(), log.Dummy)
		// Serve in background.
		go func() {
			s.Serve()
		}()

		// Create our client.
		testCli, err := tgrpc.NewTestClient(l.Addr().String())
		require.NoError(err)
		defer testCli.Close()

		// Make call.
		ns := &pbns.NodeState{
			Id:    test.id,
			State: test.state,
		}
		_, err = testCli.NodeStatusHeartbeat(context.Background(), ns)

		// Check.
		if test.shouldErr {
			assert.Error(err)
		} else {
			assert.NoError(err)
		}
		// Assert correct calls on our logic.
		mnss.AssertExpectations(t)
	}
}

func TestMasterGRPCServiceServerFailureStateList(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		name          string
		nID           *pbfs.NodeId
		fs            []*chaosv1.Failure
		expFs         []*pbfs.Failure
		stUpdateTimes int
	}{
		{
			name: "receive one failure status correctly",
			nID:  &pbfs.NodeId{Id: "test1"},
			fs: []*chaosv1.Failure{
				&chaosv1.Failure{
					Metadata: chaosv1.FailureMetadata{
						ID: "f1",
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.EnabledFailureState,
						ExpectedState: chaosv1.EnabledFailureState,
					},
				},
				&chaosv1.Failure{
					Metadata: chaosv1.FailureMetadata{
						ID: "f2",
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.EnabledFailureState,
						ExpectedState: chaosv1.EnabledFailureState,
					},
				},
				&chaosv1.Failure{
					Metadata: chaosv1.FailureMetadata{
						ID: "f3",
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.EnabledFailureState,
						ExpectedState: chaosv1.DisabledFailureState,
					},
				},
			},
			expFs: []*pbfs.Failure{
				&pbfs.Failure{Id: "f1", CurrentState: pbfs.State_ENABLED, ExpectedState: pbfs.State_ENABLED, Definition: "{}\n"},
				&pbfs.Failure{Id: "f2", CurrentState: pbfs.State_ENABLED, ExpectedState: pbfs.State_ENABLED, Definition: "{}\n"},
				&pbfs.Failure{Id: "f3", CurrentState: pbfs.State_ENABLED, ExpectedState: pbfs.State_DISABLED, Definition: "{}\n"},
			},
			stUpdateTimes: 1,
		},
		{
			name: "receive multiple failure status correctly",
			nID:  &pbfs.NodeId{Id: "test2"},
			fs: []*chaosv1.Failure{
				&chaosv1.Failure{
					Metadata: chaosv1.FailureMetadata{
						ID: "f1",
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.EnabledFailureState,
						ExpectedState: chaosv1.EnabledFailureState,
					},
				},
				&chaosv1.Failure{
					Metadata: chaosv1.FailureMetadata{
						ID: "f2",
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.EnabledFailureState,
						ExpectedState: chaosv1.EnabledFailureState,
					},
				},
				&chaosv1.Failure{
					Metadata: chaosv1.FailureMetadata{
						ID: "f3",
					},
					Status: chaosv1.FailureStatus{
						CurrentState:  chaosv1.EnabledFailureState,
						ExpectedState: chaosv1.DisabledFailureState,
					},
				},
			},
			expFs: []*pbfs.Failure{
				&pbfs.Failure{Id: "f1", CurrentState: pbfs.State_ENABLED, ExpectedState: pbfs.State_ENABLED, Definition: "{}\n"},
				&pbfs.Failure{Id: "f2", CurrentState: pbfs.State_ENABLED, ExpectedState: pbfs.State_ENABLED, Definition: "{}\n"},
				&pbfs.Failure{Id: "f3", CurrentState: pbfs.State_ENABLED, ExpectedState: pbfs.State_DISABLED, Definition: "{}\n"},
			},
			stUpdateTimes: 5,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			// Mocks.
			mnss := &mservice.NodeStatusService{}
			mfss := &mservice.FailureStatusService{}
			mfss.On("GetNodeFailures", test.nID.GetId()).Times(test.stUpdateTimes).Return(test.fs, nil)

			mclk := &mclock.Clock{}
			mclkT := make(chan time.Time)
			mclk.On("NewTicker", mock.Anything).Once().Return(&time.Ticker{C: mclkT})
			// Send the tickers N times (simulate N sends from the server).
			go func() {
				for i := 0; i < test.stUpdateTimes; i++ {
					mclkT <- time.Now()
				}
			}()

			// Create our server.
			l, err := net.Listen("tcp", "127.0.0.1:0") // :0 for a random port.
			require.NoError(err)
			defer l.Close()

			// Create our server
			s := server.NewMasterGRPCServiceServer(mfss, mnss, l, mclk, log.Dummy)

			// Serve in background.
			go func() {
				s.Serve()
			}()

			// Create our client.
			testCli, err := tgrpc.NewTestClient(l.Addr().String())
			require.NoError(err)
			defer testCli.Close()

			// Make the call.
			stream, err := testCli.FailureStatusFailureStateList(context.Background(), test.nID)

			// Check.
			if assert.NoError(err) {
				// Assert status N times (once per update).
				for i := 0; i < test.stUpdateTimes; i++ {
					fes, err := stream.Recv()
					assert.NoError(err)
					assert.Equal(test.expFs, fes.GetFailures())
				}
			}
			mfss.AssertExpectations(t)
		})
	}

}

func TestMasterGRPCServiceServerGetFailure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		name       string
		failureID  *pbfs.FailureId
		expFailure *pbfs.Failure
		expErr     bool
	}{
		{
			name:      "Correct GetFailure GRPC request",
			failureID: &pbfs.FailureId{Id: "test1"},
			expFailure: &pbfs.Failure{
				Id:            "test1",
				NodeID:        "test1node",
				Definition:    "{}\n",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: pbfs.State_DISABLED,
			},
			expErr: false,
		},
		{
			name:      "Error GetFailure GRPC request",
			failureID: &pbfs.FailureId{Id: "test2"},
			expFailure: &pbfs.Failure{
				Id:            "test2",
				NodeID:        "test2node",
				Definition:    "{}\n",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: pbfs.State_DISABLED,
			},
			expErr: true,
		},
	}

	for _, test := range tests {

		t.Run(test.name, func(t *testing.T) {
			var expErr error
			if test.expErr {
				expErr = errors.New("wanted error")
			}

			// Convert pb Failure to failure.Failure.
			cs, err := types.FailureStateTransformer.PBToFailureState(test.expFailure.GetCurrentState())
			require.NoError(err)
			es, err := types.FailureStateTransformer.PBToFailureState(test.expFailure.GetExpectedState())
			require.NoError(err)

			spec, err := chaosv1.ReadFailureSpec([]byte(test.expFailure.Definition))
			require.NoError(err)

			expF := &chaosv1.Failure{
				Metadata: chaosv1.FailureMetadata{
					ID:     test.expFailure.GetId(),
					NodeID: test.expFailure.GetNodeID(),
				},
				Spec: spec,
				Status: chaosv1.FailureStatus{
					CurrentState:  cs,
					ExpectedState: es,
				},
			}

			// Mocks.
			mnss := &mservice.NodeStatusService{}
			mfss := &mservice.FailureStatusService{}
			mfss.On("GetFailure", test.failureID.GetId()).Once().Return(expF, expErr)

			// Create our server.
			l, err := net.Listen("tcp", "127.0.0.1:0") // :0 for a random port.
			require.NoError(err)
			defer l.Close()

			// Create our server
			s := server.NewMasterGRPCServiceServer(mfss, mnss, l, clock.Base(), log.Dummy)

			// Serve in background.
			go func() {
				s.Serve()
			}()

			// Create our client.
			testCli, err := tgrpc.NewTestClient(l.Addr().String())
			require.NoError(err)
			defer testCli.Close()

			// Make the call.
			f, err := testCli.FailureStatusGetFailure(context.Background(), test.failureID)

			// Check.
			if test.expErr {
				assert.Error(err)
			} else {
				if assert.NoError(err) {
					assert.Equal(test.expFailure, f)
				}
			}
			mfss.AssertExpectations(t)
		})
	}

}
