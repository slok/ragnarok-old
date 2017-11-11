package v1_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	clusterv1pb "github.com/slok/ragnarok/api/cluster/v1/pb"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/log"
)

func TestJSONEncodeCluserV1Node(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

	tests := []struct {
		name       string
		node       *clusterv1.Node
		expEncNode string
		expErr     bool
	}{
		{
			name: "Simple object encoding shouldn't return an error if doesn't have kind or version",
			node: &clusterv1.Node{
				Metadata: clusterv1.NodeMetadata{
					ID:     "testNode1",
					Master: true,
				},
				Spec: clusterv1.NodeSpec{
					Labels: map[string]string{
						"kind": "node",
						"id":   "testNode1",
					},
				},
				Status: clusterv1.NodeStatus{
					Creation: t1,
					State:    clusterv1.ReadyNodeState,
				},
			},
			expEncNode: `{"kind":"node","version":"cluster/v1","metadata":{"id":"testNode1","master":true},"spec":{"labels":{"id":"testNode1","kind":"node"}},"status":{"state":1,"creation":"2012-11-01T22:08:41Z"}}`,
			expErr:     false,
		},
		{
			name: "Simple object encoding shouldn't return an error",
			node: &clusterv1.Node{
				TypeMeta: api.TypeMeta{
					Kind:    clusterv1.NodeKind,
					Version: clusterv1.NodeVersion,
				},
				Metadata: clusterv1.NodeMetadata{
					ID:     "testNode1",
					Master: true,
				},
				Spec: clusterv1.NodeSpec{
					Labels: map[string]string{
						"kind": "node",
						"id":   "testNode1",
					},
				},
				Status: clusterv1.NodeStatus{
					Creation: t1,
					State:    clusterv1.ReadyNodeState,
				},
			},
			expEncNode: `{"kind":"node","version":"cluster/v1","metadata":{"id":"testNode1","master":true},"spec":{"labels":{"id":"testNode1","kind":"node"}},"status":{"state":1,"creation":"2012-11-01T22:08:41Z"}}`,
			expErr:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewJSONSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
			var b bytes.Buffer
			err := s.Encode(test.node, &b)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.Equal(test.expEncNode, strings.TrimSuffix(b.String(), "\n"))
				assert.NoError(err)
			}
		})
	}
}

func TestJSONDecodeCluserV1Node(t *testing.T) {
	t1s := "2012-11-01T22:08:41Z"
	t1, _ := time.Parse(time.RFC3339, t1s)

	tests := []struct {
		name     string
		nodeJSON string
		expNode  *clusterv1.Node
		expErr   bool
	}{
		{
			name: "Simple object decoding shouldn't return an error",
			nodeJSON: `
{
	"version": "cluster/v1",
	"kind": "node",
	"metadata":{
		"id": "testNode1",
		"master": true
	},
	"spec":{
		"labels":{
			"id": "testNode1",
			"kind": "node"
		}
	},
	"status":{
		"state": 1,
		"creation": "2012-11-01T22:08:41Z"
	}
}`,
			expNode: &clusterv1.Node{
				TypeMeta: api.TypeMeta{Version: clusterv1.NodeVersion, Kind: clusterv1.NodeKind},
				Metadata: clusterv1.NodeMetadata{
					ID:     "testNode1",
					Master: true,
				},
				Spec: clusterv1.NodeSpec{
					Labels: map[string]string{
						"kind": "node",
						"id":   "testNode1",
					},
				},
				Status: clusterv1.NodeStatus{
					Creation: t1,
					State:    clusterv1.ReadyNodeState,
				},
			},
			expErr: false,
		},
		{
			name: "Simple object decoding without kind or version should return an error",
			nodeJSON: `
{
	"metadata":{
		"id": "testNode1",
		"master": true
	},
	"spec":{
		"labels":{
			"id": "testNode1",
			"kind": "node"
		}
	},
	"status":{
		"state": 1,
		"creation": "2012-11-01T22:08:41Z"
	}
}`,

			expNode: &clusterv1.Node{},
			expErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewJSONSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
			obj, err := s.Decode([]byte(test.nodeJSON))

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				node := obj.(*clusterv1.Node)
				assert.Equal(test.expNode, node)
			}
		})
	}
}

func TestYAMLEncodeCluserV1Node(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

	tests := []struct {
		name       string
		node       *clusterv1.Node
		expEncNode string
		expErr     bool
	}{
		{
			name: "Simple object encoding shouldn't return an error if doesn't have kind or version",
			node: &clusterv1.Node{
				Metadata: clusterv1.NodeMetadata{
					ID:     "testNode1",
					Master: true,
				},
				Spec: clusterv1.NodeSpec{
					Labels: map[string]string{
						"kind": "node",
						"id":   "testNode1",
					},
				},
				Status: clusterv1.NodeStatus{
					Creation: t1,
					State:    clusterv1.ReadyNodeState,
				},
			},
			expEncNode: "kind: node\nmetadata:\n  id: testNode1\n  master: true\nspec:\n  labels:\n    id: testNode1\n    kind: node\nstatus:\n  creation: 2012-11-01T22:08:41Z\n  state: 1\nversion: cluster/v1",
			expErr:     false,
		},
		{
			name: "Simple object encoding shouldn't return an error",
			node: &clusterv1.Node{
				TypeMeta: api.TypeMeta{
					Kind:    clusterv1.NodeKind,
					Version: clusterv1.NodeVersion,
				},
				Metadata: clusterv1.NodeMetadata{
					ID:     "testNode1",
					Master: true,
				},
				Spec: clusterv1.NodeSpec{
					Labels: map[string]string{
						"kind": "node",
						"id":   "testNode1",
					},
				},
				Status: clusterv1.NodeStatus{
					Creation: t1,
					State:    clusterv1.ReadyNodeState,
				},
			},
			expEncNode: "kind: node\nmetadata:\n  id: testNode1\n  master: true\nspec:\n  labels:\n    id: testNode1\n    kind: node\nstatus:\n  creation: 2012-11-01T22:08:41Z\n  state: 1\nversion: cluster/v1",
			expErr:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewYAMLSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
			var b bytes.Buffer
			err := s.Encode(test.node, &b)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.Equal(test.expEncNode, strings.TrimSuffix(b.String(), "\n"))
				assert.NoError(err)
			}
		})
	}
}

func TestYAMLDecodeCluserV1Node(t *testing.T) {
	t1s := "2012-11-01T22:08:41Z"
	t1, _ := time.Parse(time.RFC3339, t1s)

	tests := []struct {
		name     string
		nodeYAML string
		expNode  *clusterv1.Node
		expErr   bool
	}{
		{
			name: "Simple object decoding shouldn't return an error",
			nodeYAML: `
kind: node
version: cluster/v1
metadata:
  id: testNode1
  master: true
spec:
  labels:
    id: testNode1
    kind: node
status:
  creation: 2012-11-01T22:08:41Z
  state: 1`,
			expNode: &clusterv1.Node{
				TypeMeta: api.TypeMeta{Version: clusterv1.NodeVersion, Kind: clusterv1.NodeKind},
				Metadata: clusterv1.NodeMetadata{
					ID:     "testNode1",
					Master: true,
				},
				Spec: clusterv1.NodeSpec{
					Labels: map[string]string{
						"kind": "node",
						"id":   "testNode1",
					},
				},
				Status: clusterv1.NodeStatus{
					Creation: t1,
					State:    clusterv1.ReadyNodeState,
				},
			},
			expErr: false,
		},
		{
			name: "Simple object decoding without kind or version should return an error",
			nodeYAML: `
metadata:
  id: testNode1
  master: true
spec:
  labels:
    id: testNode1
    kind: node
status:
  creation: 2012-11-01T22:08:41Z
  state: 1`,
			expNode: &clusterv1.Node{},
			expErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewYAMLSerializer(serializer.ObjTyper, serializer.ObjFactory, log.Dummy)
			obj, err := s.Decode([]byte(test.nodeYAML))

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				node := obj.(*clusterv1.Node)
				assert.Equal(test.expNode, node)
			}
		})
	}
}

func TestPBEncodeCluserV1Node(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2012-11-01T22:08:41+00:00")

	tests := []struct {
		name       string
		node       *clusterv1.Node
		expEncNode *clusterv1pb.Node
		expErr     bool
	}{
		{
			name: "Simple object encoding shouldn't return an error if doesn't have kind or version",
			node: &clusterv1.Node{
				Metadata: clusterv1.NodeMetadata{
					ID:     "testNode1",
					Master: true,
				},
				Spec: clusterv1.NodeSpec{
					Labels: map[string]string{
						"kind": "node",
						"id":   "testNode1",
					},
				},
				Status: clusterv1.NodeStatus{
					Creation: t1,
					State:    clusterv1.ReadyNodeState,
				},
			},
			expEncNode: &clusterv1pb.Node{
				SerializedData: `{"kind":"node","version":"cluster/v1","metadata":{"id":"testNode1","master":true},"spec":{"labels":{"id":"testNode1","kind":"node"}},"status":{"state":1,"creation":"2012-11-01T22:08:41Z"}}`,
			},
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewPBSerializer(log.Dummy)
			pbNode := &clusterv1pb.Node{}
			err := s.Encode(test.node, pbNode)

			if test.expErr {
				assert.Error(err)
			} else {
				// Small fix for the \n
				pbNode.SerializedData = strings.TrimSuffix(pbNode.SerializedData, "\n")
				assert.Equal(test.expEncNode, pbNode)
				assert.NoError(err)
			}
		})
	}
}

func TestPBDecodeCluserV1Node(t *testing.T) {
	t1s := "2012-11-01T22:08:41Z"
	t1, _ := time.Parse(time.RFC3339, t1s)

	tests := []struct {
		name    string
		nodePB  *clusterv1pb.Node
		expNode *clusterv1.Node
		expErr  bool
	}{
		{
			name: "Simple object decoding shouldn't return an error",
			nodePB: &clusterv1pb.Node{
				SerializedData: `
{
	"version": "cluster/v1",
	"kind": "node",
	"metadata":{
		"id": "testNode1",
		"master": true
	},
	"spec":{
		"labels":{
			"id": "testNode1",
			"kind": "node"
		}
	},
	"status":{
		"state": 1,
		"creation": "2012-11-01T22:08:41Z"
	}
}`,
			},
			expNode: &clusterv1.Node{
				TypeMeta: api.TypeMeta{Version: clusterv1.NodeVersion, Kind: clusterv1.NodeKind},
				Metadata: clusterv1.NodeMetadata{
					ID:     "testNode1",
					Master: true,
				},
				Spec: clusterv1.NodeSpec{
					Labels: map[string]string{
						"kind": "node",
						"id":   "testNode1",
					},
				},
				Status: clusterv1.NodeStatus{
					Creation: t1,
					State:    clusterv1.ReadyNodeState,
				},
			},
			expErr: false,
		},
		{
			name: "Simple object decoding without kind or version should return an error",
			nodePB: &clusterv1pb.Node{
				SerializedData: `
{
	"metadata":{
		"id": "testNode1",
		"master": true
	},
	"spec":{
		"labels":{
			"id": "testNode1",
			"kind": "node"
		}
	},
	"status":{
		"state": 1,
		"creation": "2012-11-01T22:08:41Z"
	}
}`,
			},
			expNode: &clusterv1.Node{},
			expErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := serializer.NewPBSerializer(log.Dummy)
			obj, err := s.Decode(test.nodePB)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				node := obj.(*clusterv1.Node)
				assert.Equal(test.expNode, node)
			}
		})
	}
}
