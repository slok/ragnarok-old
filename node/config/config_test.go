package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/node/config"
)

func TestConfigValidation(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		masterAddr  string
		debug       bool
		dryrun      bool
		hbInterval  time.Duration
		expectError bool
	}{
		{"127.0.0.1:1234", false, true, 30 * time.Second, false},
		{"127.0.0.1:1234", true, false, 30 * time.Second, false},
		{"127.0.0.1:1234", true, false, 0 * time.Second, true},
		{"", false, true, 30 * time.Second, true},
	}

	for _, test := range tests {
		cfg := &config.Config{
			MasterAddress:     test.masterAddr,
			Debug:             test.debug,
			HeartbeatInterval: test.hbInterval,
			DryRun:            test.dryrun,
		}
		err := cfg.Validate()
		if test.expectError {
			assert.NotNil(err)
		} else {
			assert.Nil(err)
		}

	}
}
