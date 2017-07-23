package grpc_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	pb "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/model"
	"github.com/slok/ragnarok/master/service/grpc"
	mservice "github.com/slok/ragnarok/mocks/service"
	"github.com/slok/ragnarok/types"
)

func TestFailureStatusGRPCGetFailureOK(t *testing.T) {
	assert := assert.New(t)

	expF := &pb.Failure{
		Id:            "test1",
		NodeID:        "node1",
		Definition:    "definition",
		CurrentState:  pb.State_ENABLED,
		ExpectedState: pb.State_DISABLED,
	}
	stubF := &model.Failure{
		ID:            expF.GetId(),
		NodeID:        expF.GetNodeID(),
		Definition:    expF.GetDefinition(),
		CurrentState:  types.EnabledFailureState,
		ExpectedState: types.DisabledFailureState,
	}

	// Create mocks.
	mfss := &mservice.FailureStatusService{}
	mfss.On("GetFailure", expF.Id).Once().Return(stubF, nil)

	// Create the GRPC service.
	fs := grpc.NewFailureStatus(mfss, types.FailureStateTransformer, log.Dummy)

	// Get the failure and check.
	fID := &pb.FailureId{Id: stubF.ID}
	gotF, err := fs.GetFailure(context.Background(), fID)
	if assert.NoError(err) {
		assert.Equal(expF, gotF)
	}
}

func TestFailureStatusGRPCGetFailureError(t *testing.T) {
	assert := assert.New(t)

	// Create mocks.
	mfss := &mservice.FailureStatusService{}
	mfss.On("GetFailure", mock.AnythingOfType("string")).Once().Return(nil, errors.New("wanted error"))

	// Create the GRPC service.
	fs := grpc.NewFailureStatus(mfss, types.FailureStateTransformer, log.Dummy)

	// Get the failure and check.
	fID := &pb.FailureId{Id: "test"}
	_, err := fs.GetFailure(context.Background(), fID)
	assert.Error(err)
}

func TestFailureStatusGRPCGetFailureCtxCanceled(t *testing.T) {
	assert := assert.New(t)

	// Create mocks.
	mfss := &mservice.FailureStatusService{}
	mfss.On("GetFailure", mock.AnythingOfType("string")).Once().Return(nil, nil)

	// Create the GRPC service.
	fs := grpc.NewFailureStatus(mfss, types.FailureStateTransformer, log.Dummy)

	// Cancel context.
	ctx, cncl := context.WithCancel(context.Background())
	cncl()

	// Get the failure and check.
	fID := &pb.FailureId{Id: "test"}
	_, err := fs.GetFailure(ctx, fID)
	assert.Error(err)
}
