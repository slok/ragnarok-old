package client

import (
	"context"
	"fmt"
	"io"
	"sync"

	"google.golang.org/grpc"

	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	chaosv1pb "github.com/slok/ragnarok/api/chaos/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/clock"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
)

// FailureStateHandler is a custom type that knows how to handle the expected state
// of the failures.
type FailureStateHandler interface {
	// ProcessFailureStates processes a list of failures that should be in the expected state
	// this is enabled or disabled (each one will be passed in a list of that state).
	ProcessFailureStates(failures []*chaosv1.Failure) error
}

// Failure interface will implement the required methods to be able to
// communicate with a failure status server.
type Failure interface {
	// GetFailure requests and returns a Failure usinig the ID of the failure
	GetFailure(id string) (*chaosv1.Failure, error)
	// ProcessFailureStateStreaming will make a request and start reading the stream from the GRPC to handle the states.
	// It receives a handler that will be executed on every status. also receives a stop channel that will cancel the stream processing.
	ProcessFailureStateStreaming(nodeID string, handler FailureStateHandler, stopCh <-chan struct{}) error
}

// FailureGRPC staisfies Failure interface with GRPC communication.
type FailureGRPC struct {
	c          pbfs.FailureStatusClient
	serializer serializer.Serializer
	clock      clock.Clock
	logger     log.Logger

	// isStreaming will track the state when there is streaming already for a given node.
	// TODO: change to sync maps when upgrading to Go 1.9
	isStreaming     map[string]bool
	isStreamingLock sync.Mutex
}

// NewFailureGRPCFromConnection returns a new FailureGRPC using a grpc connection.
func NewFailureGRPCFromConnection(connection *grpc.ClientConn, serializer serializer.Serializer, clock clock.Clock, logger log.Logger) (*FailureGRPC, error) {
	c := pbfs.NewFailureStatusClient(connection)
	return NewFailureGRPC(c, serializer, clock, logger)
}

// NewFailureGRPC returns a new FailureGRPC.
func NewFailureGRPC(client pbfs.FailureStatusClient, serializer serializer.Serializer, clock clock.Clock, logger log.Logger) (*FailureGRPC, error) {
	return &FailureGRPC{
		c:           client,
		serializer:  serializer,
		clock:       clock,
		logger:      logger,
		isStreaming: map[string]bool{},
	}, nil
}

func (f *FailureGRPC) decodePBFailureSlice(pbFailures []*chaosv1pb.Failure) ([]*chaosv1.Failure, error) {
	res := make([]*chaosv1.Failure, len(pbFailures))
	for i, pbflr := range pbFailures {
		flr, err := f.serializer.Decode(pbflr)
		if err != nil {
			return nil, err
		}
		res[i] = flr.(*chaosv1.Failure)
	}
	return res, nil
}

// GetFailure satisfies Failure interface.
func (f *FailureGRPC) GetFailure(id string) (*chaosv1.Failure, error) {
	logger := f.logger.WithField("call", "get-failure").WithField("failureID", id)
	logger.Debug("making GRPC service call")

	// Make the call.
	fid := &pbfs.FailureId{Id: id}
	pbfl, err := f.c.GetFailure(context.Background(), fid)
	if err != nil {
		return nil, err
	}

	// transform our failure.
	flTmp, err := f.serializer.Decode(pbfl)
	if err != nil {
		return nil, fmt.Errorf("could not convert protobuf failure to internal failure type: %v", err)
	}
	fl := flTmp.(*chaosv1.Failure)
	return fl, nil
}

// ProcessFailureStateStreaming satisfies Failure interface.
func (f *FailureGRPC) ProcessFailureStateStreaming(nodeID string, handler FailureStateHandler, stopCh <-chan struct{}) error {
	logger := f.logger.WithField("call", "failure-state-list").WithField("NodeID", nodeID)
	logger.Debug("making GRPC service call")

	f.isStreamingLock.Lock()
	is, ok := f.isStreaming[nodeID]
	f.isStreamingLock.Unlock()
	if ok && is {
		return fmt.Errorf("already streaming node %s", nodeID)
	}

	// Make the call.
	nid := &pbfs.NodeId{Id: nodeID}
	stream, err := f.c.FailureStateList(context.Background(), nid)
	if err != nil {
		return err
	}

	// Start processing the stream
	f.isStreamingLock.Lock()
	f.isStreaming[nodeID] = true
	f.isStreamingLock.Unlock()

	f.logger.Info("failure status streaming started")

	go func() {
		defer func() {
			f.isStreamingLock.Lock()
			f.isStreaming[nodeID] = false
			f.isStreamingLock.Unlock()
		}()
		for {
			// Check if we have finished.
			select {
			case <-stream.Context().Done():
				f.logger.Warnf("failure status streaming terminated due context cancelation")
				if err := stream.CloseSend(); err != nil {
					f.logger.Errorf("error when closing stream: %v", err)
				}
				return
			case <-stopCh:
				f.logger.Info("failure status streaming stopped")
				return
			default:
			}

			fs, err := stream.Recv()
			pbflrs := fs.GetFailures()
			f.logger.Debugf("received: %d failures", len(pbflrs))
			if pbflrs == nil || len(pbflrs) == 0 {
				continue
			}

			if err == io.EOF {
				f.logger.Info("failure status streaming finished (server EOF)")
				return
			} else if err != nil {
				f.logger.Errorf("error receiving statuses: %v", err)
				// TODO: Reconnect
				return
			}
			transFs, err := f.decodePBFailureSlice(pbflrs)
			if err := handler.ProcessFailureStates(transFs); err != nil {
				f.logger.Errorf("error handling node %s failures: %v", nodeID, err)
			}
		}
	}()
	return nil
}
