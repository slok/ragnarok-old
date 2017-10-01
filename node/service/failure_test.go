package service_test

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/ragnarok/attack"
	"github.com/slok/ragnarok/clock"
	"github.com/slok/ragnarok/failure"
	"github.com/slok/ragnarok/log"
	mclock "github.com/slok/ragnarok/mocks/clock"
	mlog "github.com/slok/ragnarok/mocks/log"
	mclient "github.com/slok/ragnarok/mocks/node/client"
	"github.com/slok/ragnarok/node/service"
)

func TestLogFailureStateStartCli(t *testing.T) {
	tests := []struct {
		name   string
		cliErr bool
		expErr bool
	}{
		{
			name:   "If client errors when starts the handling it should error too.",
			cliErr: true,
			expErr: true,
		},
		{
			name:   "If client starts processing failures without error, the service should start handling the failure states.",
			cliErr: false,
			expErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			var cliErr error
			if test.cliErr {
				cliErr = errors.New("wanted error")
			}

			// Mocks.
			mf := &mclient.Failure{}

			lfs := service.NewLogFailureState("test", mf, clock.Base(), log.Dummy)

			// Expectetations to test.
			mf.On("ProcessFailureStateStreaming", mock.Anything, lfs, mock.Anything).Once().Return(cliErr)

			// Call logic to test and check.
			err := lfs.StartHandling()
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
			mf.AssertExpectations(t)
		})
	}
}

func TestLogFailureStateStopGoodCli(t *testing.T) {
	tests := []struct {
		name        string
		prevRunning bool
		timeout     bool
		expErr      bool
	}{
		{
			name:        "Stop correctly when no timeout and already running",
			prevRunning: true,
			timeout:     false,
			expErr:      false,
		},
		{
			name:        "If the handler is not running then it should error when trying to stop.",
			prevRunning: false,
			timeout:     false,
			expErr:      true,
		},
		{
			name:        "If the stop channel doesn't respond it will timeout and should error.",
			prevRunning: true,
			timeout:     true,
			expErr:      true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)
			require := require.New(t)

			// Mocks.
			mf := &mclient.Failure{}
			mf.On("ProcessFailureStateStreaming", mock.Anything, mock.Anything, mock.Anything).Return(nil).Run(func(args mock.Arguments) {
				// Listen to the reception of the stop signal.
				c := args.Get(2).(<-chan struct{})
				go func() {
					<-c
				}()
			})
			mc := &mclock.Clock{}
			waitC := make(chan time.Time)
			if test.timeout {
				close(waitC)
			}
			mc.On("After", mock.Anything).Return((<-chan time.Time)(waitC))

			// Create and start handling if required.
			lfs := service.NewLogFailureState("test", mf, mc, log.Dummy)
			if test.prevRunning {
				err := lfs.StartHandling()
				require.NoError(err)

			}

			err := lfs.StopHandling()
			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}
		})
	}
}

func TestLogFailureStateProcessing(t *testing.T) {
	tests := []struct {
		name     string
		failures []*failure.Failure
	}{
		{
			name:     "No failures received shouldn't log anything.",
			failures: []*failure.Failure{},
		},
		{
			name: "One failure received should log one failure.",
			failures: []*failure.Failure{
				&failure.Failure{
					ID:     "test1",
					NodeID: "node1",
					Definition: failure.Definition{
						Attacks: []failure.AttackMap{
							{
								"attack1": attack.Opts{
									"size": 524288000,
								},
							},
						},
					},
					Creation: clock.Now().UTC(),
				},
			},
		},
		{
			name: "Multiple failures received should log all the failures.",
			failures: []*failure.Failure{
				&failure.Failure{
					ID:     "test1",
					NodeID: "node1",
					Definition: failure.Definition{
						Attacks: []failure.AttackMap{
							{
								"attack1": attack.Opts{
									"size": 524288000,
								},
							},
						},
					},
					Creation: clock.Now().UTC(),
				},
				&failure.Failure{
					ID:     "test2",
					NodeID: "node1",
					Definition: failure.Definition{
						Attacks: []failure.AttackMap{},
					},
					Creation: clock.Now().UTC(),
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks.
			mf := &mclient.Failure{}
			ml := &mlog.Logger{}
			ml.On("WithField", mock.Anything, mock.Anything).Return(ml)

			// Expected calls.
			for _, expFailure := range test.failures {
				ml.On("Infof", mock.Anything, expFailure).Return()
			}

			lfs := service.NewLogFailureState("test", mf, clock.Base(), ml)
			err := lfs.ProcessFailureStates(test.failures)

			// Check has been call the logging correctly.
			if assert.NoError(err) {
				ml.AssertExpectations(t)
			}
		})
	}
}
