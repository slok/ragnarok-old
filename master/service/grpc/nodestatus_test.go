package grpc_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	pb "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service/grpc"
	mservice "github.com/slok/ragnarok/mocks/service"
	mtypes "github.com/slok/ragnarok/mocks/types"
	"github.com/slok/ragnarok/types"
)

func TestNodeStatusGRPCRegisterOK(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, types.NodeStateTransformer, log.Dummy)
	n := &pb.Node{
		Id:   "test1",
		Tags: map[string]string{"key1": "value1"},
	}

	// Mock service calls on master.
	nss.On("Register", n.Id, n.Tags).Once().Return(nil)

	// Call and check.
	resp, err := ns.Register(context.Background(), n)
	if assert.NoError(err) {
		expResp := &pb.RegisteredResponse{
			Message: fmt.Sprintf("node '%s' registered on master", n.Id),
		}
		assert.Equal(expResp, resp)
		nss.AssertExpectations(t)
	}
}

func TestNodeStatusGRPCRegisterError(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, types.NodeStateTransformer, log.Dummy)
	n := &pb.Node{
		Id:   "test1",
		Tags: map[string]string{"key1": "value1"},
	}

	// Mock service calls on master.
	nss.On("Register", n.Id, n.Tags).Once().Return(errors.New("wanted error"))

	// Call and check.
	resp, err := ns.Register(context.Background(), n)
	if assert.Error(err) {
		expResp := &pb.RegisteredResponse{
			Message: fmt.Sprintf("couldn't register node '%s' on master: %v", n.Id, err),
		}
		assert.Equal(expResp, resp)
		nss.AssertExpectations(t)
	}
}

func TestNodeStatusGRPCRegisterDoneContext(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, types.NodeStateTransformer, log.Dummy)
	n := &pb.Node{
		Id:   "test1",
		Tags: map[string]string{"key1": "value1"},
	}

	// Cancel context.
	ctx, cncl := context.WithCancel(context.Background())
	cncl()

	// Call and check.
	resp, err := ns.Register(ctx, n)

	if assert.Error(err) {
		expResp := &pb.RegisteredResponse{
			Message: "context was cancelled, not registered: context canceled",
		}
		assert.Equal(expResp, resp)
		nss.AssertExpectations(t)
	}
}

func TestNodeStatusGRPCHeartbeatOK(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}
	nsp := &mtypes.NodeStateParser{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, nsp, log.Dummy)
	n := &pb.NodeState{
		Id:    "test1",
		State: pb.State_READY,
	}

	// Mock service calls on master.
	nsp.On("PBToNodeState", mock.AnythingOfType("nodestatus.State")).Once().Return(types.ReadyNodeState, nil)
	nss.On("Heartbeat", n.Id, mock.AnythingOfType("types.NodeState")).Once().Return(nil)

	// Call and check.
	_, err := ns.Heartbeat(context.Background(), n)
	assert.NoError(err)
	nss.AssertExpectations(t)
}

func TestNodeStatusGRPCHeartbeatDoneContext(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}
	nsp := &mtypes.NodeStateParser{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, nsp, log.Dummy)

	// Cancel context.
	ctx, cncl := context.WithCancel(context.Background())
	cncl()

	// Call and check.
	_, err := ns.Heartbeat(ctx, &pb.NodeState{})
	assert.Error(err)
	nss.AssertExpectations(t)
}

func TestNodeStatusGRPCHeartbeatError(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}
	nsp := &mtypes.NodeStateParser{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, nsp, log.Dummy)

	// Mock service calls on master.
	nsp.On("PBToNodeState", mock.AnythingOfType("nodestatus.State")).Once().Return(types.ReadyNodeState, nil)
	nss.On("Heartbeat", mock.AnythingOfType("string"), types.ReadyNodeState).Once().Return(errors.New("wanted error"))

	// Call and check.
	_, err := ns.Heartbeat(context.Background(), &pb.NodeState{})
	assert.Error(err)
	nss.AssertExpectations(t)
}

func TestNodeStatusGRPCHeartbeatParseStatusError(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}
	nsp := &mtypes.NodeStateParser{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, nsp, log.Dummy)

	// Mock service calls on master.
	nsp.On("PBToNodeState", mock.AnythingOfType("nodestatus.State")).Once().Return(types.ReadyNodeState, errors.New("wanted error"))

	// Call and check.
	_, err := ns.Heartbeat(context.Background(), &pb.NodeState{})
	assert.Error(err)
	nss.AssertExpectations(t)
}
