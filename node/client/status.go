package client

import (
	"context"

	"google.golang.org/grpc"

	"github.com/slok/ragnarok/api/cluster/v1"
	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/types"
)

// Status interface will implement the required methods to be able to communicate
// with a node status server
type Status interface {
	// RegisterNode registers a node as available on the server
	RegisterNode(id string, tags map[string]string) error
	// NodeHeartbeat sends a node heartbeat to the master
	NodeHeartbeat(id string, status v1.NodeState) error
}

// StatusGRPC satisfies Status interface with GRPC communication
type StatusGRPC struct {
	c           pbns.NodeStatusClient
	stateParser types.NodeStateParser
	logger      log.Logger
}

// NewStatusGRPCFromConnection returns a new Status GRPC client based on the grpc connection
func NewStatusGRPCFromConnection(connection *grpc.ClientConn, stateParser types.NodeStateParser, logger log.Logger) (*StatusGRPC, error) {
	c := pbns.NewNodeStatusClient(connection)
	return NewStatusGRPC(c, stateParser, logger)
}

// NewStatusGRPC returns a new Status GRPC client
func NewStatusGRPC(client pbns.NodeStatusClient, stateParser types.NodeStateParser, logger log.Logger) (*StatusGRPC, error) {
	logger = logger.WithField("service", "node-status").WithField("service-kind", "grpc")
	return &StatusGRPC{
		c:           client,
		stateParser: stateParser,
		logger:      logger,
	}, nil
}

// RegisterNode satisfies Status interface
func (s *StatusGRPC) RegisterNode(id string, tags map[string]string) error {
	logger := s.logger.WithField("call", "register-node").WithField("id", id)
	logger.Debug("making GRPC service call")

	// Create the request objects
	n := &pbns.Node{
		Id:   id,
		Tags: tags,
	}

	// Make the request synchronously
	resp, err := s.c.Register(context.Background(), n)

	// Call error
	if err != nil {
		return err
	}
	// If we have a response then (we should)
	if resp != nil {
		logger.Debugf("call response: %s", resp.Message)
	}

	return nil
}

// NodeHeartbeat satisfies Status interface
func (s *StatusGRPC) NodeHeartbeat(id string, state v1.NodeState) error {
	logger := s.logger.WithField("call", "node-heartbeat").WithField("id", id)
	logger.Debug("making GRPC service call")

	st, err := s.stateParser.NodeStateToPB(state)
	if err != nil {
		return err
	}

	ns := &pbns.NodeState{
		Id:    id,
		State: st,
	}

	if _, err := s.c.Heartbeat(context.Background(), ns); err != nil {
		return err
	}
	logger.Debugf("heartbeat succeeded")
	return nil
}
