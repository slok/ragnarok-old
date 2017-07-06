package config_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/slok/ragnarok/master/config"
)

func TestConfigValidation(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		httpAddr string
		rpcAddr  string
		debug    bool

		expectError bool
	}{
		{"0.0.0.0:4444", "0.0.0.0:4441", false, false},
		{"0.0.0.0:4444", "0.0.0.0:4441", true, false},
		{"", "0.0.0.0:4441", true, true},
		{"0.0.0.0:4444", "", true, true},
		{"0.0.0.0:4444", "0.0.0.0:4444", false, true},
	}

	for _, test := range tests {
		cfg := &config.Config{
			HTTPListenAddress: test.httpAddr,
			RPCListenAddress:  test.rpcAddr,
			Debug:             test.debug,
		}
		err := cfg.Validate()
		if test.expectError {
			assert.NotNil(err)
		} else {
			assert.Nil(err)
		}

	}
}
