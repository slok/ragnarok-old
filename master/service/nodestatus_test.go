package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master"
	"github.com/slok/ragnarok/master/config"
	"github.com/slok/ragnarok/master/service"
	mmaster "github.com/slok/ragnarok/mocks/master"
	"github.com/slok/ragnarok/types"
)

func TestNodeStatusCreation(t *testing.T) {
	assert := assert.New(t)
	reg := master.NewMemNodeRepository()
	m := service.NewNodeStatus(config.Config{}, reg, log.Dummy)
	assert.NotNil(m)
}

func TestNodeStatusNodeRegistration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	n := &master.Node{
		ID:    "test1",
		Tags:  map[string]string{"address": "127.0.0.45"},
		State: types.UnknownNodeState,
	}

	// Get our registry mock.
	mReg := &mmaster.NodeRepository{}
	mReg.On("StoreNode", n.ID, n).Once().Return(nil)

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mReg, log.Dummy)
	require.NotNil(ns)

	// Check our registered node.
	err := ns.Register(n.ID, n.Tags)
	if assert.NoError(err) {
		mReg.AssertExpectations(t)
	}
}

func TestNodeStatusNodeRegistrationError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	n := &master.Node{
		ID:    "test1",
		Tags:  map[string]string{"address": "127.0.0.45"},
		State: types.UnknownNodeState,
	}

	// Get our registry mock.
	mRep := &mmaster.NodeRepository{}
	mRep.On("StoreNode", n.ID, n).Once().Return(errors.New("want error"))

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mRep, log.Dummy)
	require.NotNil(ns)

	// Check our registered node.
	err := ns.Register(n.ID, n.Tags)
	if assert.Error(err) {
		mRep.AssertExpectations(t)
	}
}

func TestNodeStatusNodeHeartbeat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	stubN := &master.Node{
		ID:    "test1",
		Tags:  map[string]string{"address": "127.0.0.45"},
		State: types.UnknownNodeState,
	}
	expN := &master.Node{
		ID:    stubN.ID,
		Tags:  stubN.Tags,
		State: types.ReadyNodeState,
	}

	// Get our repository mock.
	mRep := &mmaster.NodeRepository{}
	mRep.On("GetNode", expN.ID).Once().Return(stubN, true)
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
	mRep := &mmaster.NodeRepository{}
	mRep.On("GetNode", mock.AnythingOfType("string")).Return(nil, false)

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mRep, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat("test1", types.ReadyNodeState)
	assert.Error(err)
}

func TestFailureMasterNodeHeartbeatStoreFailure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Get our repository mock.
	mRep := &mmaster.NodeRepository{}
	mRep.On("GetNode", mock.AnythingOfType("string")).Return(&master.Node{}, true)
	mRep.On("StoreNode", mock.AnythingOfType("string"), mock.AnythingOfType("*master.Node")).Return(errors.New("wanted error"))

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mRep, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat("test1", types.ReadyNodeState)
	assert.Error(err)
}
