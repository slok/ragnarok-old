package node_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/mocks"
	mclient "github.com/slok/ragnarok/mocks/node/client"
	"github.com/slok/ragnarok/node"
	"github.com/slok/ragnarok/node/config"
	"github.com/slok/ragnarok/types"
)

func TestFailureNodeCreation(t *testing.T) {
	assert := assert.New(t)

	scm := &mclient.Status{}
	n := node.NewFailureNode(config.Config{}, scm, clock.Base(), log.Dummy)
	if assert.NotNil(n) {
		assert.NotEmpty(n.GetID())
	}
}

func TestFailureNodeCreationDryRun(t *testing.T) {
	assert := assert.New(t)

	// Mocks
	scm := &mclient.Status{}
	logger := &mocks.Logger{}
	logger.On("WithField", "id", mock.AnythingOfType("string")).Once().Return(logger)
	logger.On("Info", "System failure node ready").Once()
	logger.On("Warn", "System failure node in dry run mode").Once()

	// Check
	n := node.NewFailureNode(config.Config{DryRun: true}, scm, clock.Base(), logger)
	if assert.NotNil(n) {
		assert.NotEmpty(n.GetID())
	}
}

func TestFailureNodeRegisterOnMasterOK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create the mock
	scm := &mclient.Status{}

	// Create fNode and get the ID.
	n := node.NewFailureNode(config.Config{DryRun: true}, scm, clock.Base(), log.Dummy)
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
	scm := &mclient.Status{}

	// Create fNode and get the ID.
	cfg := config.Config{DryRun: true}
	n := node.NewFailureNode(cfg, scm, clock.Base(), log.Dummy)
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

func TestFailureNodeStartHeartbeatAlreadyRunningError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create the mock.
	msc := &mclient.Status{}
	mClock := &mocks.Clock{}
	stubT := time.NewTicker(1 * time.Millisecond)
	defer stubT.Stop()                   // Stop after test.
	heartbeatting := make(chan struct{}) // Channel that will wait for the signal when the node is already sending heartbeats.
	mClock.On("NewTicker", mock.AnythingOfType("time.Duration")).Once().Return(stubT).Run(func(args mock.Arguments) {
		heartbeatting <- struct{}{}
	})
	msc.On("NodeHeartbeat", mock.Anything, mock.Anything).Return(nil)

	// Create fNode and get the ID.
	cfg := config.Config{}
	n := node.NewFailureNode(cfg, msc, mClock, log.Dummy)
	require.NotNil(n)

	// Start first heartbeat in background.
	go func() {
		assert.NoError(n.StartHeartbeat())
	}()

	// Wait for the Goroutine started.
	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("timeout waiting for the heartbeat")
		return
	case <-heartbeatting:
		// Run & Check.
		err := n.StartHeartbeat()
		assert.Error(err)
		return
	}
}

func TestFailureNodeStartHeartbeatOK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Create the mock.
	msc := &mclient.Status{}
	mClock := &mocks.Clock{}
	stubT := time.NewTicker(1 * time.Millisecond)
	defer stubT.Stop()                    // Stop after test.
	heartbeated := make(chan struct{}, 1) // Channel that will wait for the signal when the node heartbeated.
	mClock.On("NewTicker", mock.AnythingOfType("time.Duration")).Once().Return(stubT)

	// Create fNode.
	cfg := config.Config{}
	n := node.NewFailureNode(cfg, msc, mClock, log.Dummy)
	require.NotNil(n)

	// Mock heartbeat call.
	msc.On("NodeHeartbeat", n.GetID(), types.UnknownNodeState).Return(nil).Run(func(args mock.Arguments) {
		heartbeated <- struct{}{}
	})

	// Start first heartbeat in background.
	go func() {
		assert.NoError(n.StartHeartbeat())
	}()

	// Wait for the Goroutine started.

	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("timeout waiting for the heartbeat")
		return
	case <-heartbeated:
		// Everything ok.
		msc.AssertExpectations(t)
		return
	}
}

func TestFailureNodeStoptHeartbeatError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cfg := config.Config{}

	// Create the mock.
	msc := &mclient.Status{}

	n := node.NewFailureNode(cfg, msc, clock.Base(), log.Dummy)
	require.NotNil(n)

	err := n.StopHeartbeat()
	assert.Error(err)
}

func TestFailureNodeStoptHeartbeatOK(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	cfg := config.Config{}

	// Create the mock.
	msc := &mclient.Status{}
	mClock := &mocks.Clock{}
	heartbeated := make(chan struct{})       // Channel that will wait for the signal when the node heartbeated.
	heartbeatFinished := make(chan struct{}) // Channel that will wait until heartbeat finishes.
	stubT := time.NewTicker(1 * time.Millisecond)

	// Mock calls.
	msc.On("NodeHeartbeat", mock.Anything, types.UnknownNodeState).Return(nil).Run(func(args mock.Arguments) {
		heartbeated <- struct{}{}
	})
	mClock.On("NewTicker", mock.AnythingOfType("time.Duration")).Once().Return(stubT)

	n := node.NewFailureNode(cfg, msc, mClock, log.Dummy)
	require.NotNil(n)

	// Start heartbeat in background.
	go func() {
		assert.NoError(n.StartHeartbeat())
		close(heartbeatFinished)
	}()

	// First wait unil start heartbeating.

	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("timeout waiting for the heartbeat")
		return
	case <-heartbeated:
		// Everything ok.
		msc.AssertExpectations(t)
	}

	// We are ready (already heartbeating9) to stop the heartbeating and check the stop.
	err := n.StopHeartbeat()
	require.NoError(err)

	// Check that the heartbeat stopped.
	select {
	case <-time.After(10 * time.Millisecond):
		t.Fatal("timeout waiting for the heartbeat to finish")
		return
	case <-heartbeatFinished:
		// Everything ok.
		msc.AssertExpectations(t)
		return
	}
}
