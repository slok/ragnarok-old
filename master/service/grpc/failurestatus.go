package grpc

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/slok/ragnarok/clock"
	pb "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service"
	"github.com/slok/ragnarok/types"
)

// FailureStatus implements the required methods for the FailureStatus GRPC service.
type FailureStatus struct {
	service             service.FailureStatusService // The service that has the real logic.
	fsParser            types.FailureStateParser
	stateUpdateInterval time.Duration // The interval the server will send the sate of the failures to the client.
	clock               clock.Clock
	logger              log.Logger
}

// NewFailureStatus returns a new FailureStatus.
func NewFailureStatus(stateUpdateInterval time.Duration, service service.FailureStatusService, fsParser types.FailureStateParser, clock clock.Clock, logger log.Logger) *FailureStatus {
	return &FailureStatus{
		service:             service,
		fsParser:            fsParser,
		stateUpdateInterval: stateUpdateInterval,
		clock:               clock,
		logger:              logger,
	}
}

func (f *FailureStatus) getEnabledFailureIDs(nodeID string) []string {
	fs := f.service.GetNodeExpectedEnabledFailures(nodeID)
	res := make([]string, len(fs))
	for i, flr := range fs {
		res[i] = flr.ID
	}
	return res

}
func (f *FailureStatus) getDisabledFailureIDs(nodeID string) []string {
	fs := f.service.GetNodeExpectedDisabledFailures(nodeID)
	res := make([]string, len(fs))
	for i, flr := range fs {
		res[i] = flr.ID
	}
	return res
}

// FailureStateList returns periodically the state of the current state of the failures.
func (f *FailureStatus) FailureStateList(nodeID *pb.NodeId, stream pb.FailureStatus_FailureStateListServer) error {
	f.logger.Debugf("start node %s failure update loop", nodeID.GetId())

	// Start the loop of state update for the client.
	t := f.clock.NewTicker(f.stateUpdateInterval)

	for range t.C {
		select {
		case <-stream.Context().Done():
			// Cancelled.
			f.logger.Warnf("stream update loop canceled due to context cancellation")
			return nil
		default:
		}

		// Send the state to the client.
		fes := &pb.FailuresExpectedState{
			EnabledFailureId:  f.getEnabledFailureIDs(nodeID.GetId()),
			DisabledFailureId: f.getDisabledFailureIDs(nodeID.GetId()),
		}
		if err := stream.Send(fes); err != nil {
			return fmt.Errorf("stream update loop canceled: %v", err)
		}
	}
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
