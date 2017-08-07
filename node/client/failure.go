package client

import (
	"context"
	"fmt"
	"io"
	"sync"

	"google.golang.org/grpc"

	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/failure"
	pbfs "github.com/slok/ragnarok/grpc/failurestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/types"
)

// Failure interface will implement the required methods to be able to
// communicate with a failure status server.
type Failure interface {
	// GetFailure requests and returns a Failure usinig the ID of the failure
	GetFailure(id string) (*failure.Failure, error)
	// FailureStateList will make a request and return two channels, the first one will
	// stream the failures that are in enabled state, and the second one the ones that are
	// in disabled state.
	FailureStateList(nodeID string) (enabledFs <-chan string, disabledFs <-chan string, err error)
}

// FailureGRPC staisfies Failure interface with GRPC communication.
type FailureGRPC struct {
	c             pbfs.FailureStatusClient
	stateParser   types.FailureStateParser
	failureParser failure.Parser
	clock         clock.Clock
	logger        log.Logger

	stStreaming      bool
	streamMu         sync.Mutex
	finishStreamingC chan struct{}
	// bufferLen is the number of statuses that the returning failure status channels can buffer.
	bufferLen int
}

// NewFailureGRPCFromConnection returns a new FailureGRPC using a grpc connection.
func NewFailureGRPCFromConnection(bufferLen int, connection *grpc.ClientConn, failureParser failure.Parser, stateParser types.FailureStateParser, clock clock.Clock, logger log.Logger) (*FailureGRPC, error) {
	c := pbfs.NewFailureStatusClient(connection)
	return NewFailureGRPC(bufferLen, c, failureParser, stateParser, clock, logger)
}

// NewFailureGRPC returns a new FailureGRPC.
func NewFailureGRPC(bufferLen int, client pbfs.FailureStatusClient, failureParser failure.Parser, stateParser types.FailureStateParser, clock clock.Clock, logger log.Logger) (*FailureGRPC, error) {
	return &FailureGRPC{
		c:                client,
		stateParser:      stateParser,
		failureParser:    failureParser,
		clock:            clock,
		logger:           logger,
		finishStreamingC: make(chan struct{}),
		bufferLen:        bufferLen,
	}, nil
}

// GetFailure satisfies Failure interface.
func (f *FailureGRPC) GetFailure(id string) (*failure.Failure, error) {
	logger := f.logger.WithField("call", "get-failure").WithField("failureID", id)
	logger.Debug("making GRPC service call")

	// Make the call.
	fid := &pbfs.FailureId{Id: id}
	fl, err := f.c.GetFailure(context.Background(), fid)
	if err != nil {
		return nil, err
	}

	// transform our failure.
	res, err := f.failureParser.PBToFailure(fl)
	if err != nil {
		return nil, fmt.Errorf("could not convert protobuf failure to internal failure type: %v", err)
	}

	return res, nil
}

// FailureStateList satisfies Failure interface.
func (f *FailureGRPC) FailureStateList(nodeID string) (<-chan string, <-chan string, error) {
	logger := f.logger.WithField("call", "failure-state-list").WithField("NodeID", nodeID)
	logger.Debug("making GRPC service call")

	// Before creating a stream, check if we are already streaming, if we are
	// straming the statuses then finish previous stream
	f.streamMu.Lock()
	isStreaming := f.stStreaming
	f.streamMu.Unlock()
	if isStreaming {
		f.finishStreamingC <- struct{}{}
	}

	// Make the call.
	nid := &pbfs.NodeId{Id: nodeID}
	stream, err := f.c.FailureStateList(context.Background(), nid)
	if err != nil {
		return nil, nil, err
	}

	// Make the channels.
	esC := make(chan string, f.bufferLen)
	dsC := make(chan string, f.bufferLen)

	go func() {
		// Set the streaming status to false.
		defer func() {
			f.streamMu.Lock()
			f.stStreaming = false
			f.streamMu.Unlock()
		}()

		// Start capturing the streaming
		f.logger.Info("failure status streaming started")
		f.streamMu.Lock()
		f.stStreaming = true
		f.streamMu.Unlock()

		for {
			// Check if wi have finished.
			select {
			case <-stream.Context().Done():
				f.logger.Warnf("failure status streaming terminated due con context cancelation")
				if err := stream.CloseSend(); err != nil {
					f.logger.Errorf("error when closing stream: %v", err)
				}
				return
			case <-f.finishStreamingC:
				f.logger.Info("failure status streaming finished")
				if err := stream.CloseSend(); err != nil {
					f.logger.Errorf("error when closing stream: %v", err)
				}
			default:
			}

			fs, err := stream.Recv()
			if fs == nil {
				continue
			}

			if err == io.EOF {
				return
			} else if err != nil {
				f.logger.Errorf("error receiving statuses: %v", err)
				// TODO: Reconnect
				return
			}
			// Send failures in the required channel.
			for _, fl := range fs.EnabledFailureId {
				esC <- fl
			}

			for _, fl := range fs.DisabledFailureId {
				dsC <- fl
			}
		}
	}()

	return esC, dsC, nil
}
