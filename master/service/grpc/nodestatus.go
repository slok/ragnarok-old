package grpc

import (
	"fmt"

	pbempty "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context" // TODO: Change when GRPC supports std librarie context.

	pb "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service"
	"github.com/slok/ragnarok/types"
)

// NodeStatus implements the required GRPC service methods for node status service.
type NodeStatus struct {
	service  service.NodeStatusService // The service that has the real logic
	nsParser types.NodeStateParser
	logger   log.Logger
}

// NewNodeStatus returns a new NodeStatus.
func NewNodeStatus(service service.NodeStatusService, nsParser types.NodeStateParser, logger log.Logger) *NodeStatus {
	return &NodeStatus{
		service:  service,
		nsParser: nsParser,
		logger:   logger,
	}
}

// Register registers a node on the master.
func (n *NodeStatus) Register(ctx context.Context, node *pb.Node) (*pb.RegisteredResponse, error) {
	n.logger.WithField("node", node.GetId()).Debugf("node registration GRPC call received")

	// Check context already cancelled.
	select {
	case <-ctx.Done():
		return &pb.RegisteredResponse{
			Message: fmt.Sprintf("context was cancelled, not registered: %v", ctx.Err()),
		}, ctx.Err()
	default:
	}

	if err := n.service.Register(node.Id, node.Tags); err != nil {
		return &pb.RegisteredResponse{
			Message: fmt.Sprintf("couldn't register node '%s' on master: %v", node.Id, err),
		}, err
	}

	return &pb.RegisteredResponse{
		Message: fmt.Sprintf("node '%s' registered on master", node.Id),
	}, nil
}

// Heartbeat sets the current status of a node.
func (n *NodeStatus) Heartbeat(ctx context.Context, state *pb.NodeState) (*pbempty.Empty, error) {
	n.logger.WithField("node", state.GetId()).Debugf("node heartbeat GRPC call received")

	// Check context already cancelled.
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Transform the pb state to our state kind.
	st, err := n.nsParser.PBToNodeState(state.State)
	if err != nil {
		return nil, fmt.Errorf("wrong node state: %v", state.State)
	}

	// Set the node heartbeat.
	if err := n.service.Heartbeat(state.Id, st); err != nil {
		return nil, err
	}

	return &pbempty.Empty{}, nil
}
