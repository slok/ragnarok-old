package client_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pbns "github.com/slok/ragnarok/grpc/nodestatus"
	"github.com/slok/ragnarok/log"
	mgrpc "github.com/slok/ragnarok/mocks/grpc"
	mtypes "github.com/slok/ragnarok/mocks/types"
	"github.com/slok/ragnarok/node/client"
	"github.com/slok/ragnarok/types"
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
		{"test3", map[string]string{"key1": "value1"}, false},
		{"test4", map[string]string{"key1": "value1", "key2": "value2"}, true},
	}

	for _, test := range tests {

		// Mocks
		mc := &mgrpc.NodeStatusClient{}
		var expectErr error
		if test.expectError {
			expectErr = errors.New("wanted error")
		}

		expectN := &pbns.Node{
			Id:   test.id,
			Tags: test.tags,
		}
		resp := &pbns.RegisteredResponse{Message: "called"}
		mc.On("Register", mock.Anything, expectN).Once().Return(resp, expectErr)

		// Create our client and make our service call
		s, err := client.NewStatusGRPC(mc, types.NodeStateTransformer, log.Dummy)
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

func TestNodeHeartbeat(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		id              string
		state           types.NodeState
		expState        pbns.State
		expStParseError bool
		expRespError    bool
	}{
		{"test1", types.ReadyNodeState, pbns.State_READY, false, false},
		{"test2", types.AttackingNodeState, pbns.State_ATTACKING, false, false},
		{"test3", types.ReadyNodeState, pbns.State_READY, true, false},
		{"test4", types.ReadyNodeState, pbns.State_READY, false, true},
	}

	for _, test := range tests {
		var expStParseErr error
		var expRespErr error

		if test.expStParseError {
			expStParseErr = errors.New("wanted error")
		}

		if test.expRespError {
			expRespErr = errors.New("wanted error")
		}

		// Create the mocks.
		mc := &mgrpc.NodeStatusClient{}
		mstp := &mtypes.NodeStateParser{}
		ns := &pbns.NodeState{
			Id:    test.id,
			State: test.expState,
		}
		mstp.On("NodeStateToPB", mock.Anything).Once().Return(test.expState, expStParseErr)
		// Don't call the client if there is a previous error.
		if !test.expStParseError {
			mc.On("Heartbeat", mock.Anything, ns).Once().Return(nil, expRespErr)
		}

		// Create the client.
		s, err := client.NewStatusGRPC(mc, mstp, log.Dummy)
		if assert.NoError(err) {
			err := s.NodeHeartbeat(test.id, test.state)
			if test.expRespError || test.expStParseError {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
			mstp.AssertExpectations(t)
		}
	}
}
