package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/node/config"
)

func TestConfigValidation(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		httpAddr    string
		rpcAddr     string
		debug       bool
		dryrun      bool
		expectError bool
	}{
		{"0.0.0.0:4444", "0.0.0.0:4441", false, true, false},
		{"0.0.0.0:4444", "0.0.0.0:4441", true, false, false},
		{"", "0.0.0.0:4441", true, true, true},
		{"0.0.0.0:4444", "", true, true, true},
		{"0.0.0.0:4444", "0.0.0.0:4444", false, true, true},
	}

	for _, test := range tests {
		cfg := &config.Config{
			HTTPListenAddress: test.httpAddr,
			RPCListenAddress:  test.rpcAddr,
			Debug:             test.debug,
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
