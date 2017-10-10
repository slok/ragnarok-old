package grpc_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/context"

	"time"

	"github.com/slok/ragnarok/chaos/failure"
	"github.com/slok/ragnarok/clock"
	pb "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service/grpc"
	mclock "github.com/slok/ragnarok/mocks/clock"
	mpb "github.com/slok/ragnarok/mocks/grpc/failurestatus"
	mservice "github.com/slok/ragnarok/mocks/service"
	"github.com/slok/ragnarok/types"
)

func TestFailureStatusGRPCGetFailureOK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	def := failure.Definition{}
	bs, err := def.Render()
	require.NoError(err)

	expF := &pb.Failure{
		Id:            "test1",
		NodeID:        "node1",
		Definition:    string(bs),
		CurrentState:  pb.State_ENABLED,
		ExpectedState: pb.State_DISABLED,
	}
	stubF := &failure.Failure{
		ID:            expF.GetId(),
		NodeID:        expF.GetNodeID(),
		Definition:    def,
		CurrentState:  types.EnabledFailureState,
		ExpectedState: types.DisabledFailureState,
	}

	// Create mocks.
	mfss := &mservice.FailureStatusService{}
	mfss.On("GetFailure", expF.Id).Once().Return(stubF, nil)

	// Create the GRPC service.
	fs := grpc.NewFailureStatus(0, mfss, failure.Transformer, types.FailureStateTransformer, clock.Base(), log.Dummy)

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
	fs := grpc.NewFailureStatus(0, mfss, failure.Transformer, types.FailureStateTransformer, clock.Base(), log.Dummy)

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
	fs := grpc.NewFailureStatus(0, mfss, failure.Transformer, types.FailureStateTransformer, clock.Base(), log.Dummy)

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
	fss := []*failure.Failure{
		&failure.Failure{ID: "test11", CurrentState: types.EnabledFailureState, ExpectedState: types.EnabledFailureState},
		&failure.Failure{ID: "test12", CurrentState: types.EnabledFailureState, ExpectedState: types.EnabledFailureState},
		&failure.Failure{ID: "test13", CurrentState: types.EnabledFailureState, ExpectedState: types.DisabledFailureState},
	}

	// Create mocks.
	mfss := &mservice.FailureStatusService{}
	mfss.On("GetNodeFailures", nodeID.GetId()).Times(times).Return(fss)

	mstream := &mpb.FailureStatus_FailureStateListServer{}
	mstream.On("Context").Return(context.Background())
	mstream.On("Send", mock.Anything).Return(nil)

	mtime := &mclock.Clock{}
	tC := make(chan time.Time)
	mtime.On("NewTicker", mock.Anything).Once().Return(&time.Ticker{C: tC})

	// Create the GRPC service.
	fs := grpc.NewFailureStatus(1, mfss, failure.Transformer, types.FailureStateTransformer, mtime, log.Dummy)

	// Simulate the ticker that triggers the update.
	go func() {
		for i := 0; i < times; i++ {
			tC <- time.Now()
		}
		close(tC)
	}()

	// Run the failure state refresh in background.
	err := fs.FailureStateList(nodeID, mstream)
	assert.NoError(err)

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
	fs := grpc.NewFailureStatus(1, mfss, failure.Transformer, types.FailureStateTransformer, mtime, log.Dummy)

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
	mfss.On("GetNodeFailures", nodeID.GetId()).Once().Return(nil)

	mstream := &mpb.FailureStatus_FailureStateListServer{}
	mstream.On("Context").Return(context.Background())
	mstream.On("Send", mock.Anything).Return(errors.New("wanted error"))

	mtime := &mclock.Clock{}
	tC := make(chan time.Time)
	mtime.On("NewTicker", mock.Anything).Once().Return(&time.Ticker{C: tC})

	// Create the GRPC service.
	fs := grpc.NewFailureStatus(1, mfss, failure.Transformer, types.FailureStateTransformer, mtime, log.Dummy)

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
