package server_test

import (
	"errors"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context" // TODO: Change when GRPC supports std librarie context

	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/server"
	mmaster "github.com/slok/ragnarok/mocks/master"
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
		// Create master mock.
		mm := &mmaster.Master{}
		var expErr error
		if test.shouldErr {
			expErr = errors.New("wanted error")
		}
		mm.On("RegisterNode", test.id, test.tags).Once().Return(expErr)

		// Create our server.
		l, err := net.Listen("tcp", "127.0.0.1:0") // :0 for a random port.
		require.NoError(err)
		defer l.Close()
		s := server.NewMasterGRPCServiceServer(mm, l, log.Dummy)
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
		mm.AssertExpectations(t)
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
		// Create master mock.
		mm := &mmaster.Master{}
		var expErr error
		if test.shouldErr {
			expErr = errors.New("wanted error")
		}
		mm.On("NodeHeartbeat", test.id, test.expState).Once().Return(expErr)

		// Create our server.
		l, err := net.Listen("tcp", "127.0.0.1:0") // :0 for a random port.
		require.NoError(err)
		defer l.Close()
		s := server.NewMasterGRPCServiceServer(mm, l, log.Dummy)
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
		mm.AssertExpectations(t)
	}
}
