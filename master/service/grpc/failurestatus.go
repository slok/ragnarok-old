package grpc

import (
	"fmt"
	"time"

	"golang.org/x/net/context"

	chaosv1pb "github.com/slok/ragnarok/api/chaos/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/clock"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service"
)

// FailureStatus implements the required methods for the FailureStatus GRPC service.
type FailureStatus struct {
	service             service.FailureStatusService // The service that has the real logic.
	serializer          serializer.Serializer
	stateUpdateInterval time.Duration // The interval the server will send the sate of the failures to the client.
	clock               clock.Clock
	logger              log.Logger
}

// NewFailureStatus returns a new FailureStatus.
func NewFailureStatus(stateUpdateInterval time.Duration, serializer serializer.Serializer, service service.FailureStatusService, clock clock.Clock, logger log.Logger) *FailureStatus {
	return &FailureStatus{
		service:             service,
		serializer:          serializer,
		stateUpdateInterval: stateUpdateInterval,
		clock:               clock,
		logger:              logger,
	}
}

func (f *FailureStatus) getNodeFailures(nodeID string) ([]*chaosv1pb.Failure, error) {
	fss := f.service.GetNodeFailures(nodeID)
	pbfss := make([]*chaosv1pb.Failure, len(fss))
	for i, fs := range fss {
		pbfs := &chaosv1pb.Failure{}
		if err := f.serializer.Encode(fs, pbfs); err != nil {
			return pbfss, fmt.Errorf("error while converting failure to PB failure: %s", err)
		}
		pbfss[i] = pbfs
	}
	return pbfss, nil
}

// FailureStateList returns periodically the state of the current state of the failures.
func (f *FailureStatus) FailureStateList(nodeID *pbfs.NodeId, stream pbfs.FailureStatus_FailureStateListServer) error {
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
		fs := &pbfs.FailuresState{Failures: fss}
		if err := stream.Send(fs); err != nil {
			return fmt.Errorf("stream update loop canceled: %v", err)
		}
		f.logger.WithField("targetNode", nodeID).Debugf("sent %d failures to node", len(fss))
	}
	return nil
}

// GetFailure returns a failure detail.
func (f *FailureStatus) GetFailure(ctx context.Context, fID *pbfs.FailureId) (*chaosv1pb.Failure, error) {
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

	res := &chaosv1pb.Failure{}
	if err := f.serializer.Encode(flr, res); err != nil {
		return nil, fmt.Errorf("could not make the call because of marshaling error on definition: %v", err)
	}

	return res, nil
}
