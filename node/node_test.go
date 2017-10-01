package node_test

import (
	"errors"
	"testing"

	"github.com/slok/ragnarok/log"
	mservice "github.com/slok/ragnarok/mocks/node/service"
	"github.com/slok/ragnarok/node"
	"github.com/slok/ragnarok/node/config"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
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

func TestFailureNodeStart(t *testing.T) {
	tests := []struct {
		name     string
		hbErr    bool
		fhandErr bool
		expErr   bool
	}{
		{
			name:     "If every service started correctly then it shouldn't return an error.",
			hbErr:    false,
			fhandErr: false,
			expErr:   false,
		},
		{
			name:     "If heartbeat service start fails, it should return an error.",
			hbErr:    true,
			fhandErr: false,
			expErr:   true,
		},
		{
			name:     "If failure status handler service start fails, it should return an error.",
			hbErr:    false,
			fhandErr: true,
			expErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			var hbErr, fhandErr error
			if test.hbErr {
				hbErr = errors.New("wanted error")
			}
			if test.fhandErr {
				fhandErr = errors.New("wanted error")
			}

			// Mocks.
			mfs := &mservice.FailureState{}
			if !test.hbErr { // It shouldn't reach to this point if the heartbeat didn't run ok.
				mfs.On("StartHandling").Once().Return(fhandErr)
			}
			ms := &mservice.Status{}
			ms.On("StartHeartbeat", mock.Anything).Once().Return(nil, hbErr)

			n := node.NewFailureNode("node1", config.Config{}, ms, mfs, log.Dummy)
			require.NotNil(n)
			err := n.Start()

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			ms.AssertExpectations(t)
			mfs.AssertExpectations(t)
		})
	}
}

func TestFailureNodeStop(t *testing.T) {
	tests := []struct {
		name     string
		hbErr    bool
		fhandErr bool
		expErr   bool
	}{
		{
			name:     "If every service stopped correctly then it shouldn't return an error.",
			hbErr:    false,
			fhandErr: false,
			expErr:   false,
		},
		{
			name:     "If heartbeat service stop fails, it should return an error.",
			hbErr:    true,
			fhandErr: false,
			expErr:   true,
		},
		{
			name:     "If failure status handler service stop fails, it should return an error.",
			hbErr:    false,
			fhandErr: true,
			expErr:   true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			var hbErr, fhandErr error
			if test.hbErr {
				hbErr = errors.New("wanted error")
			}
			if test.fhandErr {
				fhandErr = errors.New("wanted error")
			}

			// Mocks.
			mfs := &mservice.FailureState{}
			mfs.On("StopHandling").Once().Return(fhandErr)
			ms := &mservice.Status{}
			if !test.fhandErr { // It shouldn't reach to this point if the failure status handling stop didn't run ok.
				ms.On("StopHeartbeat", mock.Anything).Once().Return(hbErr)
			}

			n := node.NewFailureNode("node1", config.Config{}, ms, mfs, log.Dummy)
			require.NotNil(n)
			err := n.Stop()

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			ms.AssertExpectations(t)
			mfs.AssertExpectations(t)
		})
	}
}
