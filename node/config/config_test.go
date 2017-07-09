package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/node/config"
)

func TestConfigValidation(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		masterAddr  string
		debug       bool
		dryrun      bool
		expectError bool
	}{
		{"127.0.0.1:1234", false, true, false},
		{"127.0.0.1:1234", true, false, false},
		{"", false, true, true},
	}

	for _, test := range tests {
		cfg := &config.Config{
			MasterAddress: test.masterAddr,
			Debug:         test.debug,
			DryRun:        test.dryrun,
		}
		err := cfg.Validate()
		if test.expectError {
			assert.NotNil(err)
		} else {
			assert.Nil(err)
		}

	}
}
