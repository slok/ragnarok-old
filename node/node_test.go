package node_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/node"
)

func TestFailureNodeCreation(t *testing.T) {
	assert := assert.New(t)
	n := node.NewFailureNode()
	if assert.NotNil(n) {
		assert.NotEmpty(n.GetID())
	}
}
