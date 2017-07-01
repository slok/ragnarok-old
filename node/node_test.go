package node_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/mocks"
	"github.com/slok/ragnarok/node"
	"github.com/slok/ragnarok/node/config"
)

func TestFailureNodeCreation(t *testing.T) {
	assert := assert.New(t)

	logger := &mocks.Logger{}
	logger.On("WithField", "id", mock.AnythingOfType("string")).Once().Return(logger)
	logger.On("Info", "System failure node ready").Once()

	n := node.NewFailureNode(config.Config{}, logger)
	if assert.NotNil(n) {
		assert.NotEmpty(n.GetID())
	}
}

func TestFailureNodeCreationDryRun(t *testing.T) {
	assert := assert.New(t)

	logger := &mocks.Logger{}
	logger.On("WithField", "id", mock.AnythingOfType("string")).Once().Return(logger)
	logger.On("Info", "System failure node ready").Once()
	logger.On("Warn", "System failure node in dry run mode").Once()

	n := node.NewFailureNode(config.Config{DryRun: true}, logger)
	if assert.NotNil(n) {
		assert.NotEmpty(n.GetID())
	}
}
