package client_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	clusterv1 "github.com/slok/ragnarok/api/cluster/v1"
	"github.com/slok/ragnarok/apimachinery/serializer"
	"github.com/slok/ragnarok/log"
	mpbns "github.com/slok/ragnarok/mocks/grpc/nodestatus"
	"github.com/slok/ragnarok/node/client"
	testpb "github.com/slok/ragnarok/test/pb"
)

func TestRegisterNode(t *testing.T) {
	tests := []struct {
		id          string
		labels      map[string]string
		expectError bool
	}{
		{"test1", map[string]string{}, false},
		{"test2", map[string]string{}, true},
		{"test3", map[string]string{"key1": "value1"}, false},
		{"test4", map[string]string{"key1": "value1", "key2": "value2"}, true},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			assert := assert.New(t)

			// Mocks
			mc := &mpbns.NodeStatusClient{}
			var expectErr error
			if test.expectError {
				expectErr = errors.New("wanted error")
			}

			expectN := testpb.CreateLabelsPBNode(test.id, test.labels, t)
			mc.On("Register", mock.Anything, expectN).Once().Return(nil, expectErr)

			// Create our client and make our service call
			s, err := client.NewStatusGRPC(mc, serializer.PBSerializerDefault, log.Dummy)
			if assert.NoError(err) {
				// Check check result is ok
				node := &clusterv1.Node{
					Metadata: clusterv1.NodeMetadata{ID: test.id},
					Spec:     clusterv1.NodeSpec{Labels: test.labels},
				}
				err := s.RegisterNode(node)
				if test.expectError {
					assert.Error(err)
				} else {
					assert.NoError(err)
				}

			}
			// Check calls where good
			mc.AssertExpectations(t)
		})
	}
}

func TestNodeHeartbeat(t *testing.T) {
	tests := []struct {
		id           string
		expRespError bool
	}{
		{"test1", false},
		{"test2", true},
	}

	for _, test := range tests {
		t.Run(test.id, func(t *testing.T) {
			assert := assert.New(t)

			var expRespErr error
			if test.expRespError {
				expRespErr = errors.New("wanted error")
			}

			// Create the mocks.
			mc := &mpbns.NodeStatusClient{}
			n := &clusterv1.Node{
				Metadata: clusterv1.NodeMetadata{ID: test.id},
				Status:   clusterv1.NodeStatus{State: clusterv1.ReadyNodeState},
			}

			mc.On("Heartbeat", mock.Anything, mock.Anything).Once().Return(nil, expRespErr)

			// Create the client.
			s, err := client.NewStatusGRPC(mc, serializer.PBSerializerDefault, log.Dummy)
			if assert.NoError(err) {
				err := s.NodeHeartbeat(n)
				if test.expRespError {
					assert.Error(err)
				} else {
					assert.NoError(err)
				}
			}
		})
	}
}
