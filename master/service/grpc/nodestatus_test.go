package grpc_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	clusterv1pb "github.com/slok/ragnarok/api/cluster/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service/grpc"
	mserializer "github.com/slok/ragnarok/mocks/apimachinery/serializer"
	mservice "github.com/slok/ragnarok/mocks/master/service"
	testpb "github.com/slok/ragnarok/test/pb"
)

func TestNodeStatusGRPCRegisterOK(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, serializer.PBSerializerDefault, log.Dummy)
	id := "test1"
	labels := map[string]string{"key1": "value1"}
	n := testpb.CreateLabelsPBNode(id, labels, t)

	// Mock service calls on master.
	nss.On("Register", id, labels).Once().Return(nil)

	// Call and check.
	_, err := ns.Register(context.Background(), n)
	if assert.NoError(err) {
		nss.AssertExpectations(t)
	}
}

func TestNodeStatusGRPCRegisterError(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, serializer.PBSerializerDefault, log.Dummy)
	id := "test1"
	labels := map[string]string{"key1": "value1"}
	n := testpb.CreateLabelsPBNode(id, labels, t)

	// Mock service calls on master.
	nss.On("Register", id, labels).Once().Return(errors.New("wanted error"))

	// Call and check.
	_, err := ns.Register(context.Background(), n)
	if assert.Error(err) {
		nss.AssertExpectations(t)
	}
}

func TestNodeStatusGRPCRegisterDoneContext(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, serializer.PBSerializerDefault, log.Dummy)
	id := "test1"
	labels := map[string]string{"key1": "value1"}
	n := testpb.CreateLabelsPBNode(id, labels, t)

	// Cancel context.
	ctx, cncl := context.WithCancel(context.Background())
	cncl()

	// Call and check.
	_, err := ns.Register(ctx, n)

	if assert.Error(err) {
		nss.AssertExpectations(t)
	}
}

func TestNodeStatusGRPCHeartbeatOK(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, serializer.PBSerializerDefault, log.Dummy)
	id := "test1"
	state := clusterv1.ReadyNodeState
	n := testpb.CreateStatePBNode(id, state, t)

	// Mock service calls on master.
	nss.On("Heartbeat", id, state).Once().Return(nil)

	// Call and check.
	_, err := ns.Heartbeat(context.Background(), n)
	assert.NoError(err)
	nss.AssertExpectations(t)
}

func TestNodeStatusGRPCHeartbeatDoneContext(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, serializer.PBSerializerDefault, log.Dummy)

	// Cancel context.
	ctx, cncl := context.WithCancel(context.Background())
	cncl()

	// Call and check.
	_, err := ns.Heartbeat(ctx, &clusterv1pb.Node{})
	assert.Error(err)
	nss.AssertExpectations(t)
}

func TestNodeStatusGRPCHeartbeatError(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, serializer.PBSerializerDefault, log.Dummy)

	// Mock service calls on master.
	nss.On("Heartbeat", mock.Anything, mock.Anything).Once().Return(errors.New("wanted error"))

	// Call and check.
	n := testpb.CreateLabelsPBNode("test1", nil, t)
	_, err := ns.Heartbeat(context.Background(), n)
	assert.Error(err)
	nss.AssertExpectations(t)
}

func TestNodeStatusGRPCHeartbeatParseStatusError(t *testing.T) {
	assert := assert.New(t)

	// Create the mocks.
	nss := &mservice.NodeStatusService{}
	mser := &mserializer.Serializer{}

	// Create the service.
	ns := grpc.NewNodeStatus(nss, mser, log.Dummy)

	// Mock service calls on master.
	mser.On("Decode", mock.Anything).Once().Return(nil, errors.New("wanted error"))

	// Call and check.
	_, err := ns.Heartbeat(context.Background(), &clusterv1pb.Node{})
	assert.Error(err)
	nss.AssertExpectations(t)
}
