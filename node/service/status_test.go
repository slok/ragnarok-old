package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/log"
	mclock "github.com/slok/ragnarok/mocks/clock"
	mclient "github.com/slok/ragnarok/mocks/node/client"
	"github.com/slok/ragnarok/node/service"
)

func TestNodeStatusRegisterOnMaster(t *testing.T) {
	tests := []struct {
		name   string
		regErr bool
		expErr bool
		expReg bool
	}{
		{
			name:   "If client errors when the node registers it needs to error and not report registered.",
			regErr: true,
			expErr: true,
			expReg: false,
		},
		{
			name:   "When registration on node success the node should report registered.",
			regErr: false,
			expErr: false,
			expReg: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			id := "test1"
			var regErr error
			if test.regErr {
				regErr = errors.New("wanted error.")
			}
			// Create the mock
			cm := &mclock.Clock{}
			scm := &mclient.Status{}
			scm.On("RegisterNode", mock.Anything).Once().Return(regErr)

			// Create
			n := clusterv1.NewNode()
			n.Metadata.ID = id
			ns := service.NewNodeStatus(&n, scm, cm, log.Dummy)

			// Check
			err := ns.RegisterOnMaster()
			if test.expErr {
				assert.Error(err)
				assert.NotEqual(clusterv1.ReadyNodeState, ns.State())
			} else {
				assert.NoError(err)
				assert.Equal(clusterv1.ReadyNodeState, ns.State())
			}
			scm.AssertExpectations(t)
		})
	}
}

func TestNodeStatusStartHeartbeat(t *testing.T) {
	tests := []struct {
		name     string
		hbErr    bool
		reg      bool
		expErr   bool
		expHBErr bool
	}{
		{
			name:     "If nore registered the heartbeat should error.",
			reg:      false,
			hbErr:    false,
			expErr:   true,
			expHBErr: false,
		},
		{
			name:     "If error on heartbeat it should return an error on the hb error channel.",
			reg:      true,
			hbErr:    true,
			expErr:   false,
			expHBErr: true,
		},
		{
			name:     "If registered and heartbeat ok shouldn't return errors.",
			reg:      true,
			hbErr:    false,
			expErr:   false,
			expHBErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			id := "test1"
			var hbErr error
			if test.hbErr {
				hbErr = errors.New("wanted error")
			}

			// Create the mock
			cm := &mclock.Clock{}
			cm.On("NewTicker", mock.Anything).Return(clock.NewTicker(1))
			cm.On("After", mock.Anything).Return(clock.After(9999999999)) // No timeout on error senders.
			scm := &mclient.Status{}
			scm.On("RegisterNode", mock.Anything).Once().Return(nil)
			hbCall := make(chan struct{})
			scm.On("NodeHeartbeat", mock.Anything).Return(hbErr).Run(func(_ mock.Arguments) {
				hbCall <- struct{}{}
			})

			// Create
			n := clusterv1.NewNode()
			n.Metadata.ID = id
			ns := service.NewNodeStatus(&n, scm, cm, log.Dummy)
			if test.reg {
				err := ns.RegisterOnMaster()
				require.NoError(err)
			}
			hbErrC, err := ns.StartHeartbeat(1)

			if test.expErr {
				assert.Error(err)
			} else {
				require.NoError(err)
				select {
				case <-clock.After(10 * time.Millisecond):
					assert.Fail("timeout waiting for heartbeat.")
				case <-hbCall:
					// Expect error on the hb errors channel.
					if test.expHBErr {
						select {
						case <-clock.After(100 * time.Millisecond):
							assert.Fail("timeout waiting for heartbeat error.")
						case err := <-hbErrC:
							assert.Error(err)
						}
					}
					// All good.
				}
			}
		})
	}
}

func TestNodeStatusStopHeartbeat(t *testing.T) {
	tests := []struct {
		name    string
		startHB bool
		expErr  bool
	}{
		{
			name:    "Stopping the heartbeat before starting it should error.",
			startHB: false,
			expErr:  true,
		},
		{
			name:    "Stopping the heartbeat after starting a heartbeat should stop the heartbeat and shouldn't return an error.",
			startHB: true,
			expErr:  false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)
			id := "test1"

			cm := &mclock.Clock{}
			cm.On("NewTicker", mock.Anything).Return(clock.NewTicker(1))
			cm.On("After", mock.Anything).Return(clock.After(9999999999)) // No timeout on error senders.
			scm := &mclient.Status{}
			scm.On("RegisterNode", mock.Anything).Once().Return(nil)
			readyC := make(chan struct{})
			scm.On("NodeHeartbeat", mock.Anything).Return(nil).Run(func(_ mock.Arguments) {
				readyC <- struct{}{}
			})

			n := clusterv1.NewNode()
			n.Metadata.ID = id
			ns := service.NewNodeStatus(&n, scm, cm, log.Dummy)
			err := ns.RegisterOnMaster()
			require.NoError(err)
			if test.startHB {
				_, err = ns.StartHeartbeat(1)
				require.NoError(err)

				// Wait until the heartbeat has started to continue the test.
				select {
				case <-clock.After(10 * time.Millisecond):
					require.Fail("timeout waiting to the heartbeat start.")
				case <-readyC:
				}
			}

			err = ns.StopHeartbeat()
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})

	}
}
