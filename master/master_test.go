package master_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master"
	"github.com/slok/ragnarok/master/config"
	mmaster "github.com/slok/ragnarok/mocks/master"
)

func TestFailureMasterCreation(t *testing.T) {
	assert := assert.New(t)
	reg := master.NewMemNodeRegistry()
	m := master.NewFailureMaster(config.Config{}, reg, log.Dummy)
	assert.NotNil(m)
}

func TestFailureMasterNodeRegistration(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	n := &master.Node{
		ID:      "test1",
		Tags: map[string]string{"address": "127.0.0.45"},
	}

	// Get our registry mock
	mReg := &mmaster.NodeRegistry{}
	mReg.On("AddNode", n.ID, n).Once().Return(nil)

	// Create a master
	m := master.NewFailureMaster(config.Config{}, mReg, log.Dummy)
	require.NotNil(n)

	// Check our registered node
	err := m.RegisterNode(n.ID, n.Tags)
	if assert.NoError(err) {
		mReg.AssertExpectations(t)
	}
}

func TestFailureMasterNodeRegistrationError(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	n := &master.Node{
		ID:   "test1",
		Tags: map[string]string{"address": "127.0.0.45"},
	}

	// Get our registry mock
	mReg := &mmaster.NodeRegistry{}
	mReg.On("AddNode", n.ID, n).Once().Return(errors.New("want error"))

	// Create a master
	m := master.NewFailureMaster(config.Config{}, mReg, log.Dummy)
	require.NotNil(n)

	// Check our registered node
	err := m.RegisterNode(n.ID, n.Tags)
	if assert.Error(err) {
		mReg.AssertExpectations(t)
	}
}
