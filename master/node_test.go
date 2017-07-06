package master_test

import (
	"testing"

	"fmt"

	"github.com/slok/ragnarok/master"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMemNodeRegistryRegisterNode(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		quantity int
	}{
		{quantity: 1},
		{quantity: 5},
		{quantity: 100},
	}

	for _, test := range tests {
		reg := master.NewMemNodeRegistry()
		nodes := make([]*master.Node, test.quantity)

		for i := 0; i < test.quantity; i++ {
			n := &master.Node{
				ID:      fmt.Sprintf("id-%d", i),
				Address: fmt.Sprintf("127.0.0.%d", i),
			}
			nodes = append(nodes, n)
			err := reg.AddNode(n.ID, n)
			assert.NoError(err)

			// Check stored node is ok
			nGot, ok := reg.GetNode(n.ID)
			if assert.True(ok) {
				assert.Equal(n, nGot)
			}
		}
	}
}

func TestMemNodeRegistryGetMissin(t *testing.T) {
	assert := assert.New(t)

	reg := master.NewMemNodeRegistry()
	nGot, ok := reg.GetNode("missing")
	if assert.False(ok) {
		assert.Nil(nGot)
	}
}
func TestMemNodeRegistryDelete(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	reg := master.NewMemNodeRegistry()

	n := &master.Node{
		ID:      "test1",
		Address: "127.0.0.1:2314",
	}
	err := reg.AddNode(n.ID, n)
	require.NoError(err)
	_, ok := reg.GetNode(n.ID)
	require.True(ok)

	// Check delete works
	reg.DeleteNode(n.ID)
	_, ok = reg.GetNode(n.ID)
	assert.False(ok)
}

func TestMemNodeRegistryAddGetAll(t *testing.T) {
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
		reg := master.NewMemNodeRegistry()
		nodes := make([]*master.Node, test.quantity)

		for i := 0; i < test.quantity; i++ {
			n := &master.Node{
				ID:      fmt.Sprintf("id-%d", i),
				Address: fmt.Sprintf("127.0.0.%d", i),
			}
			nodes = append(nodes, n)
			err := reg.AddNode(n.ID, n)
			require.NoError(err)
		}

		// Check number of nodes
		nsGot := reg.GetNodes()
		assert.Len(nsGot, test.quantity)
	}
}
