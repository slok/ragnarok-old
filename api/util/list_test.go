package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	apiutil "github.com/slok/ragnarok/api/util"
)

func TestNewObjectList(t *testing.T) {
	assert := assert.New(t)

	nodeList := []*clusterv1.Node{
		&clusterv1.Node{
			Metadata: api.ObjectMeta{
				ID: "testNode1",
				Labels: map[string]string{
					"kind": "node",
					"id":   "testNode1",
				},
				Annotations: map[string]string{
					"name": "my node",
				},
			},
			Spec: clusterv1.NodeSpec{},
			Status: clusterv1.NodeStatus{
				State: clusterv1.ReadyNodeState,
			},
		},
		&clusterv1.Node{
			Metadata: api.ObjectMeta{
				ID: "testNode2",
				Labels: map[string]string{
					"kind": "node",
					"id":   "testNode2",
				},
				Annotations: map[string]string{
					"name": "my node number 2",
				},
			},
			Spec: clusterv1.NodeSpec{},
			Status: clusterv1.NodeStatus{
				State: clusterv1.ReadyNodeState,
			},
		},
	}
	cnt := "123456"
	expObjectList := clusterv1.NewNodeList(nodeList, cnt)

	objectList := make([]api.Object, len(nodeList))
	for i, n := range nodeList {
		objectList[i] = n
	}

	gotObjectList, err := apiutil.NewObjectList(objectList, cnt)

	if assert.NoError(err) {
		assert.Equal(&expObjectList, gotObjectList)
	}

}
