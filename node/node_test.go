package node_test

import (
	"testing"

	"github.com/slok/ragnarok/log"
	mservice "github.com/slok/ragnarok/mocks/node/service"
	"github.com/slok/ragnarok/node"
	"github.com/slok/ragnarok/node/config"

	"github.com/stretchr/testify/assert"
)

func TestFailureNodeCreation(t *testing.T) {
	assert := assert.New(t)
	mfs := &mservice.FailureState{}
	ms := &mservice.Status{}
	n := node.NewFailureNode("node1", config.Config{}, ms, mfs, log.Dummy)
	if assert.NotNil(n) {
		assert.NotEmpty(n.GetID())
	}
}
