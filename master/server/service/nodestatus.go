package service

import (
	"fmt"

	"golang.org/x/net/context" // TODO: Change when GRPC supports std librarie context

	pb "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master"
)

// NodeStatusGRPC implements the required GRPC service methods for node status service.
type NodeStatusGRPC struct {
	master master.Master
	logger log.Logger
}

// NewNodeStatusGRPC returns a new NodeStatusGRPC.
func NewNodeStatusGRPC(master master.Master, logger log.Logger) *NodeStatusGRPC {
	return &NodeStatusGRPC{
		master: master,
		logger: logger,
	}
}

// Register registers a node on the master.
func (n *NodeStatusGRPC) Register(ctx context.Context, node *pb.Node) (*pb.RegisteredResponse, error) {
	// Check context already cancelled.
	select {
	case <-ctx.Done():
		return &pb.RegisteredResponse{
			Message: fmt.Sprintf("context was cancelled, not registered: %v", ctx.Err()),
		}, ctx.Err()
	default:
	}

	if err := n.master.RegisterNode(node.Id, node.Tags); err != nil {
		return &pb.RegisteredResponse{
			Message: fmt.Sprintf("couldn't register node '%s' on master: %v", node.Id, err),
		}, err
	}

	return &pb.RegisteredResponse{
		Message: fmt.Sprintf("node '%s' registered on master", node.Id),
	}, nil
}
