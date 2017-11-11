package client

import (
	"context"

	"google.golang.org/grpc"

	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	clusterv1pb "github.com/slok/ragnarok/api/cluster/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
)

// Status interface will implement the required methods to be able to communicate
// with a node status server
type Status interface {
	// RegisterNode registers a node as available on the server
	RegisterNode(node *clusterv1.Node) error
	// NodeHeartbeat sends a node heartbeat to the master
	NodeHeartbeat(node *clusterv1.Node) error
}

// StatusGRPC satisfies Status interface with GRPC communication
type StatusGRPC struct {
	c          pbns.NodeStatusClient
	serializer serializer.Serializer // pb serializer
	logger     log.Logger
}

// NewStatusGRPCFromConnection returns a new Status GRPC client based on the grpc connection
func NewStatusGRPCFromConnection(connection *grpc.ClientConn, serializer serializer.Serializer, logger log.Logger) (*StatusGRPC, error) {
	c := pbns.NewNodeStatusClient(connection)
	return NewStatusGRPC(c, serializer, logger)
}

// NewStatusGRPC returns a new Status GRPC client
func NewStatusGRPC(client pbns.NodeStatusClient, serializer serializer.Serializer, logger log.Logger) (*StatusGRPC, error) {
	logger = logger.WithField("service", "node-status").WithField("service-kind", "grpc")
	return &StatusGRPC{
		c:          client,
		serializer: serializer,
		logger:     logger,
	}, nil
}

// RegisterNode satisfies Status interface
func (s *StatusGRPC) RegisterNode(node *clusterv1.Node) error {
	logger := s.logger.WithField("call", "register-node").WithField("id", node.Metadata.ID)
	logger.Debug("making GRPC service call")

	// Create the request objects
	pbn := &clusterv1pb.Node{}
	if err := s.serializer.Encode(node, pbn); err != nil {
		return err
	}

	// Make the request synchronously
	_, err := s.c.Register(context.Background(), pbn)

	return err
}

// NodeHeartbeat satisfies Status interface
func (s *StatusGRPC) NodeHeartbeat(node *clusterv1.Node) error {
	logger := s.logger.WithField("call", "node-heartbeat").WithField("id", node.Metadata.ID)

	logger.Debug("making GRPC service call")

	// Create the request objects
	pbn := &clusterv1pb.Node{}
	if err := s.serializer.Encode(node, pbn); err != nil {
		return err
	}

	if _, err := s.c.Heartbeat(context.Background(), pbn); err != nil {
		return err
	}
	logger.Debugf("heartbeat succeeded")
	return nil
}
