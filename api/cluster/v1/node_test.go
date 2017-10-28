package v1_test

import (
	"bytes"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/apimachinery"
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
			name: "Simple object encoding should return an error if doesn't have kind or version",
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
			expEncNode: "",
			expErr:     true,
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
			expEncNode: "{\"kind\":\"node\",\"version\":\"cluster/v1\",\"metadata\":{\"id\":\"testNode1\",\"master\":true},\"spec\":{\"labels\":{\"id\":\"testNode1\",\"kind\":\"node\"}},\"status\":{\"state\":1,\"creation\":\"2012-11-01T22:08:41Z\"}}\n",
			expErr:     false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := apimachinery.NewJSONSerializer(apimachinery.ObjFactory, log.Dummy)
			var b bytes.Buffer
			err := s.Encode(test.node, &b)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.Equal(test.expEncNode, b.String())
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
			expErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			s := apimachinery.NewJSONSerializer(apimachinery.ObjFactory, log.Dummy)
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
