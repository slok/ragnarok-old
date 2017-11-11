package grpc

import (
	emptypb "github.com/golang/protobuf/ptypes/empty"
	"golang.org/x/net/context" // TODO: Change when GRPC supports std librarie context.

	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	clusterv1pb "github.com/slok/ragnarok/api/cluster/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service"
)

// NodeStatus implements the required GRPC service methods for node status service.
type NodeStatus struct {
	service    service.NodeStatusService // The service that has the real logic
	serializer serializer.Serializer
	logger     log.Logger
}

// NewNodeStatus returns a new NodeStatus.
func NewNodeStatus(service service.NodeStatusService, serializer serializer.Serializer, logger log.Logger) *NodeStatus {
	return &NodeStatus{
		service:    service,
		serializer: serializer,
		logger:     logger,
	}
}

// Register registers a node on the master.
func (n *NodeStatus) Register(ctx context.Context, nodepb *clusterv1pb.Node) (*emptypb.Empty, error) {
	empty := &emptypb.Empty{}

	// Decode pb object.
	nodeObj, err := n.serializer.Decode(nodepb)
	if err != nil {
		return empty, err
	}
	node := nodeObj.(*clusterv1.Node)

	n.logger.WithField("node", node.Metadata.ID).Debugf("node registration GRPC call received")
	// Check context already cancelled.
	select {
	case <-ctx.Done():
		return empty, ctx.Err()
	default:
	}

	if err := n.service.Register(node.Metadata.ID, node.Spec.Labels); err != nil {
		return empty, err
	}

	return empty, nil
}

// Heartbeat sets the current status of a node.
func (n *NodeStatus) Heartbeat(ctx context.Context, nodepb *clusterv1pb.Node) (*emptypb.Empty, error) {
	empty := &emptypb.Empty{}

	// Decode pb object.
	nodeObj, err := n.serializer.Decode(nodepb)
	if err != nil {
		return empty, err
	}
	node := nodeObj.(*clusterv1.Node)

	n.logger.WithField("node", node.Metadata.ID).Debugf("node heartbeat GRPC call received")
	// Check context already cancelled.
	select {
	case <-ctx.Done():
		return empty, ctx.Err()
	default:
	}

	// Set the node heartbeat.
	if err := n.service.Heartbeat(node.Metadata.ID, node.Status.State); err != nil {
		return nil, err
	}

	return empty, nil
}
