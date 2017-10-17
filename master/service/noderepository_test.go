package service_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/master/service"
)

func TestMemNodeRepositoryRegisterNode(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		quantity int
	}{
		{quantity: 1},
		{quantity: 5},
		{quantity: 100},
	}

	for _, test := range tests {
		reg := service.NewMemNodeRepository()

		for i := 0; i < test.quantity; i++ {
			n := v1.Node{
				ID:     fmt.Sprintf("id-%d", i),
				Labels: map[string]string{"address": fmt.Sprintf("127.0.0.%d", i)},
				State:  v1.ReadyNodeState,
			}
			err := reg.StoreNode(n.ID, n)
			assert.NoError(err)

			// Check stored node is ok
			nGot, ok := reg.GetNode(n.ID)
			if assert.True(ok) {
				assert.Equal(&n, nGot)
			}
		}
	}
}

func TestMemNodeRepositoryGetMissing(t *testing.T) {
	assert := assert.New(t)

	reg := service.NewMemNodeRepository()
	nGot, ok := reg.GetNode("missing")
	if assert.False(ok) {
		assert.Nil(nGot)
	}
}
func TestMemNodeRepositoryDelete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	reg := service.NewMemNodeRepository()

	n := v1.Node{
		ID:     "test1",
		Labels: map[string]string{"address": "127.0.0.1"},
		State:  v1.AttackingNodeState,
	}
	err := reg.StoreNode(n.ID, n)
	require.NoError(err)
	_, ok := reg.GetNode(n.ID)
	require.True(ok)

	// Check delete works
	reg.DeleteNode(n.ID)
	_, ok = reg.GetNode(n.ID)
	assert.False(ok)
}

func TestMemNodeRepositoryStoreGetAll(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	tests := []struct {
		quantity int
	}{
		{quantity: 1},
		{quantity: 5},
		{quantity: 100},
	}

	for _, test := range tests {
		reg := service.NewMemNodeRepository()
		nodes := make([]v1.Node, test.quantity)

		for i := 0; i < test.quantity; i++ {
			n := v1.Node{
				ID:     fmt.Sprintf("id-%d", i),
				Labels: map[string]string{"address": fmt.Sprintf("127.0.0.%d", i)},
				State:  v1.ErroredNodeState,
			}
			nodes = append(nodes, n)
			err := reg.StoreNode(n.ID, n)
			require.NoError(err)
		}

		// Check number of nodes
		nsGot := reg.GetNodes()
		assert.Len(nsGot, test.quantity)
	}
}

func TestMemNodeRepositoryGetNodesByLabels(t *testing.T) {
	tests := []struct {
		name     string
		nodes    []v1.Node
		selector v1.NodeLabels
		expNodes map[string]*v1.Node
	}{
		{
			name: "No labels shouldn't return any node",
			nodes: []v1.Node{
				v1.Node{
					ID: "node1",
					Labels: map[string]string{
						"id":   "node1",
						"env":  "prod",
						"kind": "master",
					},
				},
				v1.Node{
					ID: "node2",
					Labels: map[string]string{
						"id":   "node2",
						"env":  "staging",
						"kind": "master",
					},
				},
				v1.Node{
					ID: "node3",
					Labels: map[string]string{
						"id":   "node3",
						"env":  "prod",
						"kind": "node",
					},
				},
			},
			selector: map[string]string{},
			expNodes: map[string]*v1.Node{},
		},
		{
			name: "Single ID label should return one node only",
			nodes: []v1.Node{
				v1.Node{
					ID: "node1",
					Labels: map[string]string{
						"id":   "node1",
						"env":  "prod",
						"kind": "master",
					},
				},
				v1.Node{
					ID: "node2",
					Labels: map[string]string{
						"id":   "node2",
						"env":  "staging",
						"kind": "master",
					},
				},
				v1.Node{
					ID: "node3",
					Labels: map[string]string{
						"id":   "node3",
						"env":  "prod",
						"kind": "node",
					},
				},
			},
			selector: map[string]string{"id": "node2"},
			expNodes: map[string]*v1.Node{
				"node2": &v1.Node{
					ID: "node2",
					Labels: map[string]string{
						"id":   "node2",
						"env":  "staging",
						"kind": "master",
					},
				},
			},
		},
		{
			name: "Single ID label should return one node only",
			nodes: []v1.Node{
				v1.Node{
					ID: "node1",
					Labels: map[string]string{
						"id":   "node1",
						"env":  "prod",
						"kind": "master",
					},
				},
				v1.Node{
					ID: "node2",
					Labels: map[string]string{
						"id":   "node2",
						"env":  "staging",
						"kind": "master",
					},
				},
				v1.Node{
					ID: "node3",
					Labels: map[string]string{
						"id":   "node3",
						"env":  "prod",
						"kind": "node",
					},
				},
				v1.Node{
					ID: "node4",
					Labels: map[string]string{
						"id":   "node4",
						"env":  "prod",
						"kind": "master",
					},
				},
			},
			selector: map[string]string{"env": "prod", "kind": "master"},
			expNodes: map[string]*v1.Node{
				"node1": &v1.Node{
					ID: "node1",
					Labels: map[string]string{
						"id":   "node1",
						"env":  "prod",
						"kind": "master",
					},
				},
				"node4": &v1.Node{
					ID: "node4",
					Labels: map[string]string{
						"id":   "node4",
						"env":  "prod",
						"kind": "master",
					},
				},
			},
		},
	}

	for _, test := range tests {
		assert := assert.New(t)
		require := require.New(t)

		t.Run(test.name, func(t *testing.T) {
			reg := service.NewMemNodeRepository()

			// Insert the nodes.
			for _, n := range test.nodes {
				require.NoError(reg.StoreNode(n.ID, n))
			}

			gotN := reg.GetNodesByLabels(test.selector)

			assert.Equal(test.expNodes, gotN)
		})
	}
}
