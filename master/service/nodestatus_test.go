package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
	"github.com/slok/ragnarok/master/service"
	mservice "github.com/slok/ragnarok/mocks/service"
)

func TestNodeStatusCreation(t *testing.T) {
	assert := assert.New(t)
	reg := service.NewMemNodeRepository()
	m := service.NewNodeStatus(config.Config{}, reg, log.Dummy)
	assert.NotNil(m)
}

func TestNodeStatusNodeRegistration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	n := v1.Node{
		ID:     "test1",
		Labels: map[string]string{"address": "127.0.0.45"},
		State:  v1.UnknownNodeState,
	}

	// Get our registry mock.
	mReg := &mservice.NodeRepository{}
	mReg.On("StoreNode", n.ID, n).Once().Return(nil)

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mReg, log.Dummy)
	require.NotNil(ns)

	// Check our registered node.
	err := ns.Register(n.ID, n.Labels)
	if assert.NoError(err) {
		mReg.AssertExpectations(t)
	}
}

func TestNodeStatusNodeRegistrationError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	n := v1.Node{
		ID:     "test1",
		Labels: map[string]string{"address": "127.0.0.45"},
		State:  v1.UnknownNodeState,
	}

	// Get our registry mock.
	mRep := &mservice.NodeRepository{}
	mRep.On("StoreNode", n.ID, n).Once().Return(errors.New("want error"))

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mRep, log.Dummy)
	require.NotNil(ns)

	// Check our registered node.
	err := ns.Register(n.ID, n.Labels)
	if assert.Error(err) {
		mRep.AssertExpectations(t)
	}
}

func TestNodeStatusNodeHeartbeat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	stubN := v1.Node{
		ID:     "test1",
		Labels: map[string]string{"address": "127.0.0.45"},
		State:  v1.UnknownNodeState,
	}
	expN := v1.Node{
		ID:     stubN.ID,
		Labels: stubN.Labels,
		State:  v1.ReadyNodeState,
	}

	// Get our repository mock.
	mRep := &mservice.NodeRepository{}
	mRep.On("GetNode", expN.ID).Once().Return(&stubN, true)
	mRep.On("StoreNode", expN.ID, expN).Once().Return(nil)

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mRep, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat(expN.ID, expN.State)
	if assert.NoError(err) {
		mRep.AssertExpectations(t)
	}
}

func TestNodeStatusNodeHeartbeatNotRegistered(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Get our repository mock.
	mRep := &mservice.NodeRepository{}
	mRep.On("GetNode", mock.AnythingOfType("string")).Return(nil, false)

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mRep, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat("test1", v1.ReadyNodeState)
	assert.Error(err)
}

func TestNodeStatusNodeHeartbeatStoreFailure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Get our repository mock.
	mRep := &mservice.NodeRepository{}
	mRep.On("GetNode", mock.AnythingOfType("string")).Return(&v1.Node{}, true)
	mRep.On("StoreNode", mock.AnythingOfType("string"), mock.Anything).Return(errors.New("wanted error"))

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mRep, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat("test1", v1.ReadyNodeState)
	assert.Error(err)
}
