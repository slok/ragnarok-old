package service_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/api"
	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/config"
	"github.com/slok/ragnarok/master/service"
	mcliclusterv1 "github.com/slok/ragnarok/mocks/client/api/cluster/v1"
)

func TestNodeStatusNodeRegistration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	n := &clusterv1.Node{
		TypeMeta: api.TypeMeta{Kind: clusterv1.NodeKind, Version: clusterv1.NodeVersion},
		Metadata: api.ObjectMeta{
			ID:     "test1",
			Labels: map[string]string{"address": "127.0.0.45"},
		},
		Status: clusterv1.NodeStatus{
			State: clusterv1.UnknownNodeState,
		},
	}

	// Get our registry mock.
	mcli := &mcliclusterv1.NodeClientInterface{}
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

	n := &clusterv1.Node{
		TypeMeta: api.TypeMeta{Kind: clusterv1.NodeKind, Version: clusterv1.NodeVersion},
		Metadata: api.ObjectMeta{
			ID:     "test1",
			Labels: map[string]string{"address": "127.0.0.45"},
		},
		Status: clusterv1.NodeStatus{
			State: clusterv1.UnknownNodeState,
		},
	}

	// Get our registry mock.
	mcli := &mcliclusterv1.NodeClientInterface{}
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

	stubN := clusterv1.Node{
		TypeMeta: api.TypeMeta{Kind: clusterv1.NodeKind, Version: clusterv1.NodeVersion},
		Metadata: api.ObjectMeta{
			ID:     "test1",
			Labels: map[string]string{"address": "127.0.0.45"},
		},
		Status: clusterv1.NodeStatus{
			State: clusterv1.UnknownNodeState,
		},
	}
	expN := &clusterv1.Node{
		TypeMeta: api.TypeMeta{Kind: clusterv1.NodeKind, Version: clusterv1.NodeVersion},
		Metadata: api.ObjectMeta{
			ID:     stubN.Metadata.ID,
			Labels: stubN.Metadata.Labels,
		},
		Status: clusterv1.NodeStatus{
			State: clusterv1.ReadyNodeState,
		},
	}

	// Get our repository mock.
	mcli := &mcliclusterv1.NodeClientInterface{}
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
	mcli := &mcliclusterv1.NodeClientInterface{}
	mcli.On("Get", mock.Anything).Return(nil, fmt.Errorf("wanted error"))

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mcli, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat("test1", clusterv1.ReadyNodeState)
	assert.Error(err)
}

func TestNodeStatusNodeHeartbeatStoreFailure(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	// Get our repository mock.
	mcli := &mcliclusterv1.NodeClientInterface{}
	mcli.On("Get", mock.Anything).Return(&clusterv1.Node{}, nil)
	mcli.On("Update", mock.Anything, mock.Anything).Return(nil, errors.New("wanted error"))

	// Create the service.
	ns := service.NewNodeStatus(config.Config{}, mcli, log.Dummy)
	require.NotNil(ns)

	// Check our heartbeat node
	err := ns.Heartbeat("test1", clusterv1.ReadyNodeState)
	assert.Error(err)
}
