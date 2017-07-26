package grpc_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/net/context"

	"time"

	"github.com/slok/ragnarok/clock"
	pb "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/model"
	"github.com/slok/ragnarok/master/service/grpc"
	mclock "github.com/slok/ragnarok/mocks/clock"
	mpb "github.com/slok/ragnarok/mocks/grpc/failurestatus"
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
	fs := grpc.NewFailureStatus(0, mfss, types.FailureStateTransformer, clock.Base(), log.Dummy)

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
	fs := grpc.NewFailureStatus(0, mfss, types.FailureStateTransformer, clock.Base(), log.Dummy)

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
	fs := grpc.NewFailureStatus(0, mfss, types.FailureStateTransformer, clock.Base(), log.Dummy)

	// Cancel context.
	ctx, cncl := context.WithCancel(context.Background())
	cncl()

	// Get the failure and check.
	fID := &pb.FailureId{Id: "test"}
	_, err := fs.GetFailure(ctx, fID)
	assert.Error(err)
}

func TestFailureStatusGRPCFailureStateListOK(t *testing.T) {
	assert := assert.New(t)

	nodeID := &pb.NodeId{Id: "test1"}
	times := 5
	efs := []*model.Failure{
		&model.Failure{ID: "test11"},
		&model.Failure{ID: "test12"},
	}
	dfs := []*model.Failure{
		&model.Failure{ID: "test13"},
	}
	expectedSt := &pb.FailuresExpectedState{
		EnabledFailureId:  []string{"test11", "test12"},
		DisabledFailureId: []string{"test13"},
	}

	// Create mocks.
	mfss := &mservice.FailureStatusService{}
	mfss.On("GetNodeExpectedEnabledFailures", nodeID.GetId()).Times(times).Return(efs)
	mfss.On("GetNodeExpectedDisabledFailures", nodeID.GetId()).Times(times).Return(dfs)

	mstream := &mpb.FailureStatus_FailureStateListServer{}
	mstream.On("Context").Return(context.Background())
	mstream.On("Send", expectedSt).Return(nil)

	mtime := &mclock.Clock{}
	tC := make(chan time.Time)
	mtime.On("NewTicker", mock.Anything).Once().Return(&time.Ticker{C: tC})

	// Create the GRPC service.
	fs := grpc.NewFailureStatus(1, mfss, types.FailureStateTransformer, mtime, log.Dummy)

	// Run the failure state refresh in background.
	go func() {
		err := fs.FailureStateList(nodeID, mstream)
		assert.NoError(err)
	}()

	// Simulate the ticker that triggers the update.
	for i := 0; i < times; i++ {
		tC <- time.Now()
	}
	close(tC)

	time.Sleep(5 * time.Millisecond) // Used to wait for the final calls and have a real assert.
	mfss.AssertExpectations(t)
	mstream.AssertExpectations(t)
}

func TestFailureStatusGRPCFailureStateListContextClosed(t *testing.T) {
	assert := assert.New(t)

	nodeID := &pb.NodeId{Id: "test1"}

	// Create mocks.
	mfss := &mservice.FailureStatusService{}

	mstream := &mpb.FailureStatus_FailureStateListServer{}
	ctx, clfn := context.WithCancel(context.Background())
	clfn() // Cancel the context.
	mstream.On("Context").Return(ctx)

	mtime := &mclock.Clock{}
	tC := make(chan time.Time)
	mtime.On("NewTicker", mock.Anything).Once().Return(&time.Ticker{C: tC})

	// Create the GRPC service.
	fs := grpc.NewFailureStatus(1, mfss, types.FailureStateTransformer, mtime, log.Dummy)

	// Trigger one round on the update loop.
	go func() {
		tC <- time.Now()
	}()

	// Check.
	err := fs.FailureStateList(nodeID, mstream)
	if assert.NoError(err) {
		mfss.AssertExpectations(t)
		mstream.AssertExpectations(t)
	}

	close(tC)
}

func TestFailureStatusGRPCFailureStateListErr(t *testing.T) {
	assert := assert.New(t)

	nodeID := &pb.NodeId{Id: "test1"}

	// Create mocks.
	mfss := &mservice.FailureStatusService{}
	mfss.On("GetNodeExpectedEnabledFailures", nodeID.GetId()).Once().Return(nil)
	mfss.On("GetNodeExpectedDisabledFailures", nodeID.GetId()).Once().Return(nil)

	mstream := &mpb.FailureStatus_FailureStateListServer{}
	mstream.On("Context").Return(context.Background())
	mstream.On("Send", mock.Anything).Return(errors.New("wanted error"))

	mtime := &mclock.Clock{}
	tC := make(chan time.Time)
	mtime.On("NewTicker", mock.Anything).Once().Return(&time.Ticker{C: tC})

	// Create the GRPC service.
	fs := grpc.NewFailureStatus(1, mfss, types.FailureStateTransformer, mtime, log.Dummy)

	// Trigger one round on the update loop.
	go func() {
		tC <- time.Now()
	}()

	// Check.
	err := fs.FailureStateList(nodeID, mstream)
	if assert.Error(err) {
		mfss.AssertExpectations(t)
		mstream.AssertExpectations(t)

	}

	close(tC)
}
