package master_test

import (
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/master"
	"github.com/slok/ragnarok/types"
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
		reg := master.NewMemNodeRepository()
		nodes := make([]*master.Node, test.quantity)

		for i := 0; i < test.quantity; i++ {
			n := &master.Node{
				ID:    fmt.Sprintf("id-%d", i),
				Tags:  map[string]string{"address": fmt.Sprintf("127.0.0.%d", i)},
				State: types.ReadyNodeState,
			}
			nodes = append(nodes, n)
			err := reg.StoreNode(n.ID, n)
			assert.NoError(err)

			// Check stored node is ok
			nGot, ok := reg.GetNode(n.ID)
			if assert.True(ok) {
				assert.Equal(n, nGot)
			}
		}
	}
}

func TestMemNodeRepositoryGetMissing(t *testing.T) {
	assert := assert.New(t)

	reg := master.NewMemNodeRepository()
	nGot, ok := reg.GetNode("missing")
	if assert.False(ok) {
		assert.Nil(nGot)
	}
}
func TestMemNodeRepositoryDelete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	reg := master.NewMemNodeRepository()

	n := &master.Node{
		ID:    "test1",
		Tags:  map[string]string{"address": "127.0.0.1"},
		State: types.AttackingNodeState,
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
		reg := master.NewMemNodeRepository()
		nodes := make([]*master.Node, test.quantity)

		for i := 0; i < test.quantity; i++ {
			n := &master.Node{
				ID:    fmt.Sprintf("id-%d", i),
				Tags:  map[string]string{"address": fmt.Sprintf("127.0.0.%d", i)},
				State: types.ErroredNodeState,
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
