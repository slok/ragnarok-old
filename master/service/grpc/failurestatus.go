package grpc

import (
	"context"

	pb "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service"
	"github.com/slok/ragnarok/types"
)

// FailureStatus implements the required methods for the FailureStatus GRPC service.
type FailureStatus struct {
	service  service.FailureStatusService // The service that has the real logic.
	fsParser types.FailureStateParser
	logger   log.Logger
}

// NewFailureStatus returns a new FailureStatus.
func NewFailureStatus(service service.FailureStatusService, fsParser types.FailureStateParser, logger log.Logger) *FailureStatus {
	return &FailureStatus{
		service:  service,
		fsParser: fsParser,
		logger:   logger,
	}
}

// FailureStateList returns periodically the state of the current state of the failures.
func (f *FailureStatus) FailureStateList(*pb.NodeId, pb.FailureStatus_FailureStateListServer) error {
	return nil
}

// GetFailure returns a failure detail.
func (f *FailureStatus) GetFailure(ctx context.Context, fID *pb.FailureId) (*pb.Failure, error) {
	// Check context already cancelled.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	flr, err := f.service.GetFailure(fID.GetId())
	if err != nil {
		return nil, err
	}

	cSt, err := f.fsParser.FailureStateToPB(flr.CurrentState)
	if err != nil {
		return nil, err
	}
	eSt, err := f.fsParser.FailureStateToPB(flr.ExpectedState)
	if err != nil {
		return nil, err
	}

	res := &pb.Failure{
		Id:            flr.ID,
		NodeID:        flr.NodeID,
		Definition:    flr.Definition,
		CurrentState:  cSt,
		ExpectedState: eSt,
	}

	return res, nil
}
