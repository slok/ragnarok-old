package client_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	mgrpc "github.com/slok/ragnarok/mocks/grpc"
	"github.com/slok/ragnarok/node/client"
)

func TestRegisterNode(t *testing.T) {
	assert := assert.New(t)
	tests := []struct {
		id          string
		tags        map[string]string
		expectError bool
	}{
		{"test1", map[string]string{}, false},
		{"test2", map[string]string{}, true},
		{"test2", map[string]string{"key1": "value1"}, false},
		{"test2", map[string]string{"key1": "value1", "key2": "value2"}, true},
	}

	for _, test := range tests {

		// Mocks
		mc := &mgrpc.NodeStatusClient{}
		var expectErr error
		if test.expectError {
			expectErr = errors.New("wanted error")
		}

		expectN := &pb.Node{
			Id:   test.id,
			Tags: test.tags,
		}
		resp := &pb.RegisteredResponse{Message: "called"}
		mc.On("Register", mock.Anything, expectN).Once().Return(resp, expectErr)

		// Create our client and make our service call
		s, err := client.NewStatusGRPC(mc, log.Dummy)
		if assert.NoError(err) {
			// Check check result is ok
			err := s.RegisterNode(test.id, test.tags)
			if test.expectError {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

		}
		// Check calls where good
		mc.AssertExpectations(t)
	}
}
