package node_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/mocks"
	"github.com/slok/ragnarok/node"
	"github.com/slok/ragnarok/node/config"
)

func TestFailureNodeCreation(t *testing.T) {
	assert := assert.New(t)

	scm := &mocks.Status{}
	n := node.NewFailureNode(config.Config{}, scm, log.Dummy)
	if assert.NotNil(n) {
		assert.NotEmpty(n.GetID())
	}
}

func TestFailureNodeCreationDryRun(t *testing.T) {
	assert := assert.New(t)

	// Mocks
	scm := &mocks.Status{}
	logger := &mocks.Logger{}
	logger.On("WithField", "id", mock.AnythingOfType("string")).Once().Return(logger)
	logger.On("Info", "System failure node ready").Once()
	logger.On("Warn", "System failure node in dry run mode").Once()

	// Check
	n := node.NewFailureNode(config.Config{DryRun: true}, scm, logger)
	if assert.NotNil(n) {
		assert.NotEmpty(n.GetID())
	}
}

func TestFailureNodeRegisterOnMasterOK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create the mock
	scm := &mocks.Status{}

	// Create fNode and get the ID
	n := node.NewFailureNode(config.Config{DryRun: true}, scm, log.Dummy)
	require.NotNil(n)
	id := n.GetID()

	// Mock the call
	scm.On("RegisterNode", id, mock.AnythingOfType("map[string]string")).Once().Return(nil)

	// Check
	err := n.RegisterOnMaster()
	if assert.NoError(err) {
		scm.AssertExpectations(t)
	}
}

func TestFailureNodeRegisterOnMasterError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create the mock
	scm := &mocks.Status{}

	// Create fNode and get the ID
	n := node.NewFailureNode(config.Config{DryRun: true}, scm, log.Dummy)
	require.NotNil(n)
	id := n.GetID()

	// Mock the call
	scm.On("RegisterNode", id, mock.AnythingOfType("map[string]string")).Once().Return(errors.New(""))

	// Check
	err := n.RegisterOnMaster()
	if assert.Error(err) {
		scm.AssertExpectations(t)
	}
}
