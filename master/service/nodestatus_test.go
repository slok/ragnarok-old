package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/api/cluster/v1"
	cliclusterv1 "github.com/slok/ragnarok/client/cluster/v1"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
	"github.com/slok/ragnarok/master/service"
	mcliclusterv1 "github.com/slok/ragnarok/mocks/client/cluster/v1"
)

func TestNodeStatusCreation(t *testing.T) {
	assert := assert.New(t)
	client := cliclusterv1.NewDefaultNodeMem()
	m := service.NewNodeStatus(config.Config{}, client, log.Dummy)
	assert.NotNil(m)
}

func TestNodeStatusNodeRegistration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	n := &v1.Node{
		Metadata: api.ObjectMeta{
			ID:     "test1",
			Labels: map[string]string{"address": "127.0.0.45"},
		},
		Status: v1.NodeStatus{
			State: v1.UnknownNodeState,
		},
	}

	// Get our registry mock.
	mcli := &mcliclusterv1.Node{}
	mcli.On("Create", n).Once().Return(nil, nil)

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mcli, log.Dummy)
	require.NotNil(ns)

	// Check our registered node.
	err := ns.Register(n.Metadata.ID, n.Metadata.Labels)
	if assert.NoError(err) {
		mcli.AssertExpectations(t)
	}
}

func TestNodeStatusNodeRegistrationError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	n := &v1.Node{
		Metadata: api.ObjectMeta{
			ID:     "test1",
			Labels: map[string]string{"address": "127.0.0.45"},
		},
		Status: v1.NodeStatus{
			State: v1.UnknownNodeState,
		},
	}

	// Get our registry mock.
	mcli := &mcliclusterv1.Node{}
	mcli.On("Create", n).Once().Return(nil, errors.New("want error"))

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mcli, log.Dummy)
	require.NotNil(ns)

	// Check our registered node.
	err := ns.Register(n.Metadata.ID, n.Metadata.Labels)
	if assert.Error(err) {
		mcli.AssertExpectations(t)
	}
}

func TestNodeStatusNodeHeartbeat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	stubN := v1.Node{
		Metadata: api.ObjectMeta{
			ID:     "test1",
			Labels: map[string]string{"address": "127.0.0.45"},
		},
		Status: v1.NodeStatus{
			State: v1.UnknownNodeState,
		},
	}
	expN := &v1.Node{
		Metadata: api.ObjectMeta{
			ID:     stubN.Metadata.ID,
			Labels: stubN.Metadata.Labels,
		},
		Status: v1.NodeStatus{
			State: v1.ReadyNodeState,
		},
	}

	// Get our repository mock.
	mcli := &mcliclusterv1.Node{}
	mcli.On("Get", expN.Metadata.ID).Once().Return(&stubN, nil)
	mcli.On("Update", expN).Once().Return(nil, nil)

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mcli, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat(expN.Metadata.ID, expN.Status.State)
	if assert.NoError(err) {
		mcli.AssertExpectations(t)
	}
}

func TestNodeStatusNodeHeartbeatNotRegistered(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Get our repository mock.
	mcli := &mcliclusterv1.Node{}
	mcli.On("Get", mock.Anything).Return(nil, fmt.Errorf("wanted error"))

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mcli, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat("test1", v1.ReadyNodeState)
	assert.Error(err)
}

func TestNodeStatusNodeHeartbeatStoreFailure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Get our repository mock.
	mcli := &mcliclusterv1.Node{}
	mcli.On("Get", mock.Anything).Return(&v1.Node{}, nil)
	mcli.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("wanted error"))

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mcli, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat("test1", v1.ReadyNodeState)
	assert.Error(err)
}
