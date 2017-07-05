package client

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
)

// Status interface will implement the required methods to be able to communicate
// with a node status server
type Status interface {
	// RegisterNode registers a node as available on the server
	RegisterNode(id string, tags map[string]string) error
}

// StatusGRPC satisfies Status interface with GRPC communication
type StatusGRPC struct {
	c      pb.NodeStatusClient
	logger log.Logger
}

// NewStatusGRPCFromConnection returns a new Status GRPC client based on the grpc connection
func NewStatusGRPCFromConnection(connection *grpc.ClientConn, logger log.Logger) (*StatusGRPC, error) {
	c := pb.NewNodeStatusClient(connection)
	return NewStatusGRPC(c, logger)
}

// NewStatusGRPC returns a new Status GRPC client
func NewStatusGRPC(client pb.NodeStatusClient, logger log.Logger) (*StatusGRPC, error) {
	logger = logger.WithField("service", "node-status").WithField("service-kind", "grpc")
	return &StatusGRPC{
		c:      client,
		logger: logger,
	}, nil
}

// RegisterNode satisfies Status interface
func (s *StatusGRPC) RegisterNode(id string, tags map[string]string) error {
	logger := s.logger.WithField("call", "register-node").WithField("id", id)
	logger.Debug("making GRPC service call")

	// Create the request objects
	ni := &pb.NodeInfo{
		Node: &pb.Node{
			Id: id,
		},
		Address: "0.0.0.0", // TODO: Set up correct address
		Tags:    tags,
	}

	// Make the request synchronously
	resp, err := s.c.Register(context.Background(), ni)

	// Call error
	if err != nil {
		return err
	}
	// If we have a response then (we should, just checking)
	if resp != nil {
		logger.Debug("call response: %s", resp.Message)
	}

	return nil
}
