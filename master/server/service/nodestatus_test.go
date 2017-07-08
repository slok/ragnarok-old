package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"

	pb "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/server/service"
	mmaster "github.com/slok/ragnarok/mocks/master"
)

func TestNodeStatusGRPCRegisterOK(t *testing.T) {
	assert := assert.New(t)
	// Create the mocks.
	mm := &mmaster.Master{}

	// Create the service.
	nsg := service.NewNodeStatusGRPC(mm, log.Dummy)
	n := &pb.NodeInfo{
		Node:    &pb.Node{Id: "test1"},
		Address: "127.0.0.1:1234",
		Tags:    map[string]string{"key1": "value1"},
	}

	// Mock service calls on master.
	mm.On("RegisterNode", n.Node.Id, "127.0.0.1:1234").Once().Return(nil)

	// Call and check.
	resp, err := nsg.Register(context.Background(), n)
	if assert.NoError(err) {
		expResp := &pb.RegisteredResponse{
			Message: fmt.Sprintf("node '%s' registered on master", n.Node.Id),
		}
		assert.Equal(expResp, resp)
		mm.AssertExpectations(t)

	}
}

func TestNodeStatusGRPCRegisterError(t *testing.T) {
	assert := assert.New(t)
	// Create the mocks.
	mm := &mmaster.Master{}

	// Create the service.
	nsg := service.NewNodeStatusGRPC(mm, log.Dummy)
	n := &pb.NodeInfo{
		Node:    &pb.Node{Id: "test1"},
		Address: "127.0.0.1:1234",
		Tags:    map[string]string{"key1": "value1"},
	}

	// Mock service calls on master.
	mm.On("RegisterNode", n.Node.Id, "127.0.0.1:1234").Once().Return(errors.New("wanted error"))

	// Call and check.
	resp, err := nsg.Register(context.Background(), n)
	if assert.Error(err) {
		expResp := &pb.RegisteredResponse{
			Message: fmt.Sprintf("couldn't register node '%s' on master: %v", n.Node.Id, err),
		}
		assert.Equal(expResp, resp)
		mm.AssertExpectations(t)
	}
}

func TestNodeStatusGRPCRegisterDoneContext(t *testing.T) {
	assert := assert.New(t)
	// Create the mocks.
	mm := &mmaster.Master{}

	// Create the service.
	nsg := service.NewNodeStatusGRPC(mm, log.Dummy)
	n := &pb.NodeInfo{
		Node:    &pb.Node{Id: "test1"},
		Address: "127.0.0.1:1234",
		Tags:    map[string]string{"key1": "value1"},
	}

	// Cancel context.
	ctx, cncl := context.WithCancel(context.Background())
	cncl()

	// Call and check.
	resp, err := nsg.Register(ctx, n)

	if assert.Error(err) {
		expResp := &pb.RegisteredResponse{
			Message: "context was cancelled, not registered: context canceled",
		}
		assert.Equal(expResp, resp)
		mm.AssertExpectations(t)
	}
}
