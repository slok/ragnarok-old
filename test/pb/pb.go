package pb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	clusterv1pb "github.com/slok/ragnarok/api/cluster/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
)

// CreateLabelsPBNode helper function to create pb nodes.
func CreateLabelsPBNode(id string, labels map[string]string, t *testing.T) *clusterv1pb.Node {
	n := &clusterv1.Node{
		Metadata: clusterv1.NodeMetadata{
			ID: id,
		},
		Spec: clusterv1.NodeSpec{
			Labels: labels,
		},
	}
	return CreatePBNode(n, t)
}

// CreateStatePBNode helper function to create pb nodes.
func CreateStatePBNode(id string, state clusterv1.NodeState, t *testing.T) *clusterv1pb.Node {
	n := &clusterv1.Node{
		Metadata: clusterv1.NodeMetadata{
			ID: id,
		},
		Status: clusterv1.NodeStatus{
			State: state,
		},
	}
	return CreatePBNode(n, t)
}

// CreatePBNode helper function to create pb nodes.
func CreatePBNode(n *clusterv1.Node, t *testing.T) *clusterv1pb.Node {
	var b bytes.Buffer
	assert.NoError(t, serializer.DefaultSerializer.Encode(n, &b))

	return &clusterv1pb.Node{
		SerializedData: b.String(),
	}
}
