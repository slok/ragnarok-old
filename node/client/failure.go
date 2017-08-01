package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/slok/ragnarok/attack"
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
}

// FailureGRPC staisfies Failure interface with GRPC communication.
type FailureGRPC struct {
	c             pbfs.FailureStatusClient
	stateParser   types.FailureStateParser
	failureParser failure.Parser
	attackReg     attack.Registry // The registry where all the attacks are registered
	clock         clock.Clock
	logger        log.Logger
}

// NewFailureGRPCFromConnection returns a new FailureGRPC using a grpc connection.
func NewFailureGRPCFromConnection(connection *grpc.ClientConn, failureParser failure.Parser, stateParser types.FailureStateParser, clock clock.Clock, logger log.Logger) (*FailureGRPC, error) {
	c := pbfs.NewFailureStatusClient(connection)
	return NewFailureGRPC(c, failureParser, stateParser, clock, logger)
}

// NewFailureGRPC returns a new FailureGRPC.
func NewFailureGRPC(client pbfs.FailureStatusClient, failureParser failure.Parser, stateParser types.FailureStateParser, clock clock.Clock, logger log.Logger) (*FailureGRPC, error) {
	return &FailureGRPC{
		c:             client,
		stateParser:   stateParser,
		failureParser: failureParser,
		clock:         clock,
		logger:        logger,
	}, nil
}

// GetFailure satisfies FAilure interface.
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
