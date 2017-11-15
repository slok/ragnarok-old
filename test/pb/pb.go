package pb

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api"
	chaosv1 "github.com/slok/ragnarok/api/chaos/v1"
	chaosv1pb "github.com/slok/ragnarok/api/chaos/v1/pb"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	clusterv1pb "github.com/slok/ragnarok/api/cluster/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
)

// CreateLabelsPBNode helper function to create pb nodes.
func CreateLabelsPBNode(id string, labels map[string]string, t *testing.T) *clusterv1pb.Node {
	n := &clusterv1.Node{
		Metadata: api.ObjectMeta{
			ID:     id,
			Labels: labels,
		},
	}
	return CreatePBNode(n, t)
}

// CreateStatePBNode helper function to create pb nodes.
func CreateStatePBNode(id string, state clusterv1.NodeState, t *testing.T) *clusterv1pb.Node {
	n := &clusterv1.Node{
		Metadata: api.ObjectMeta{
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

// CreatePBFailure helper function to create pb failures.
func CreatePBFailure(f *chaosv1.Failure, t *testing.T) *chaosv1pb.Failure {
	var b bytes.Buffer
	assert.NoError(t, serializer.DefaultSerializer.Encode(f, &b))

	return &chaosv1pb.Failure{
		SerializedData: b.String(),
	}
}
