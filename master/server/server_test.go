package server_test

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context" // TODO: Change when GRPC supports std librarie context

	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/model"
	"github.com/slok/ragnarok/master/server"
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
		s := server.NewMasterGRPCServiceServer(mfss, mnss, l, log.Dummy)
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
		expState  types.NodeState
		shouldErr bool
	}{
		{"test1", pbns.State_READY, types.ReadyNodeState, false},
		{"test1", pbns.State_UNKNOWN, types.UnknownNodeState, false},
		{"test1", pbns.State_ERRORED, types.ErroredNodeState, false},
		{"test1", pbns.State_ATTACKING, types.AttackingNodeState, false},
		{"test1", pbns.State_REVERTING, types.RevertingNodeState, false},
		{"test1", pbns.State_REVERTING, types.RevertingNodeState, true},
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
		s := server.NewMasterGRPCServiceServer(mfss, mnss, l, log.Dummy)
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

func _TestMasterGRPCServiceServerFailureStateList(t *testing.T) {
	assert := assert.New(t)
	assert.Fail("Need to be implemented")
}

func TestMasterGRPCServiceServerGetFailure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		failureID  *pbfs.FailureId
		expFailure *pbfs.Failure
		expErr     bool
	}{
		{
			failureID: &pbfs.FailureId{Id: "test1"},
			expFailure: &pbfs.Failure{
				Id:            "test1",
				NodeID:        "test1node",
				Definition:    "test1definition",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: pbfs.State_DISABLED,
			},
			expErr: false,
		},
		{
			failureID: &pbfs.FailureId{Id: "test2"},
			expFailure: &pbfs.Failure{
				Id:            "test2",
				NodeID:        "test2node",
				Definition:    "test2definition",
				CurrentState:  pbfs.State_ENABLED,
				ExpectedState: pbfs.State_DISABLED,
			},
			expErr: true,
		},
	}

	for _, test := range tests {

		var expErr error
		if test.expErr {
			expErr = errors.New("wanted error")
		}

		// Converti pb Failure to model.Failure.
		cs, err := types.FailureStateTransformer.PBToFailureState(test.expFailure.GetCurrentState())
		require.NoError(err)
		es, err := types.FailureStateTransformer.PBToFailureState(test.expFailure.GetExpectedState())
		require.NoError(err)
		expF := &model.Failure{
			ID:            test.expFailure.GetId(),
			NodeID:        test.expFailure.GetNodeID(),
			Definition:    test.expFailure.GetDefinition(),
			CurrentState:  cs,
			ExpectedState: es,
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
		s := server.NewMasterGRPCServiceServer(mfss, mnss, l, log.Dummy)

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
	}

}
