package grpc

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/failure"
	pb "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service"
	"github.com/slok/ragnarok/types"
)

// FailureStatus implements the required methods for the FailureStatus GRPC service.
type FailureStatus struct {
	service             service.FailureStatusService // The service that has the real logic.
	fParser             failure.Parser
	fsParser            types.FailureStateParser
	stateUpdateInterval time.Duration // The interval the server will send the sate of the failures to the client.
	clock               clock.Clock
	logger              log.Logger
}

// NewFailureStatus returns a new FailureStatus.
func NewFailureStatus(stateUpdateInterval time.Duration, service service.FailureStatusService, fParser failure.Parser, fsParser types.FailureStateParser, clock clock.Clock, logger log.Logger) *FailureStatus {
	return &FailureStatus{
		service:             service,
		fParser:             fParser,
		fsParser:            fsParser,
		stateUpdateInterval: stateUpdateInterval,
		clock:               clock,
		logger:              logger,
	}
}

func (f *FailureStatus) getNodeFailures(nodeID string) ([]*pb.Failure, error) {
	fss := f.service.GetNodeFailures(nodeID)
	pbFSs := make([]*pb.Failure, len(fss))
	for i, fs := range fss {
		pbf, err := f.fParser.FailureToPB(fs)
		if err != nil {
			return pbFSs, fmt.Errorf("error while converting failure to PB failure: %s", err)
		}
		pbFSs[i] = pbf
	}
	return pbFSs, nil
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
		fss, err := f.getNodeFailures(nodeID.GetId())
		if err != nil {
			return err
		}
		fs := &pb.FailuresState{Failures: fss}
		if err := stream.Send(fs); err != nil {
			return fmt.Errorf("stream update loop canceled: %v", err)
		}
		f.logger.WithField("targetNode", nodeID).Debugf("sent %d failures to node", len(fss))
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

	// Marshal the definition to bytearray.
	bs, err := flr.Definition.Render()
	if err != nil {
		return nil, fmt.Errorf("could not make the call because of marshaling error on definition: %v", err)
	}

	res := &pb.Failure{
		Id:            flr.ID,
		NodeID:        flr.NodeID,
		Definition:    string(bs),
		CurrentState:  cSt,
		ExpectedState: eSt,
	}

	return res, nil
}
