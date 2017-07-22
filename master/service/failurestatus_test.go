package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/model"
	"github.com/slok/ragnarok/master/service"
	mservice "github.com/slok/ragnarok/mocks/service"
	"github.com/slok/ragnarok/types"
)

type testNodeFailures map[string][]*model.Failure

func TestGetNodeFailures(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		expectedFailures testNodeFailures
	}{
		{
			expectedFailures: testNodeFailures{},
		},
		{
			expectedFailures: testNodeFailures{
				"node1": {
					&model.Failure{ID: "f11", NodeID: "node1"},
				},
				"node2": {
					&model.Failure{ID: "f21", NodeID: "node2"},
					&model.Failure{ID: "f21", NodeID: "node2"},
				},
				"node3": {},
			},
		},
	}

	for _, test := range tests {
		// Create mocks.
		mrepo := &mservice.FailureRepository{}

		// Mock the call.
		for nID, expfs := range test.expectedFailures {
			mrepo.On("GetAllByNode", nID).Once().Return(expfs)
		}

		// Create the service.
		fs := service.NewFailureStatus(mrepo, log.Dummy)

		// Loop on every node.
		for nID, expfs := range test.expectedFailures {
			gotFs := fs.GetNodeFailures(nID)
			assert.Equal(expfs, gotFs)
		}
	}
}

func TestGetNodeExpectedEnabledFailures(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		failures    []*model.Failure
		expFailures []*model.Failure
	}{
		{
			failures: []*model.Failure{
				&model.Failure{ID: "f1", ExpectedState: types.DisabledFailureState},
			},
			expFailures: []*model.Failure{},
		},
		{
			failures: []*model.Failure{
				&model.Failure{ID: "f1", ExpectedState: types.EnabledFailureState},
			},
			expFailures: []*model.Failure{
				&model.Failure{ID: "f1", ExpectedState: types.EnabledFailureState},
			},
		},
		{
			failures: []*model.Failure{
				&model.Failure{ID: "f1", ExpectedState: types.DisabledFailureState},
				&model.Failure{ID: "f2", ExpectedState: types.EnabledFailureState},
				&model.Failure{ID: "f3", ExpectedState: types.RevertingFailureState},
				&model.Failure{ID: "f4", ExpectedState: types.EnabledFailureState},
			},
			expFailures: []*model.Failure{
				&model.Failure{ID: "f2", ExpectedState: types.EnabledFailureState},
				&model.Failure{ID: "f4", ExpectedState: types.EnabledFailureState},
			},
		},
	}

	for _, test := range tests {
		// Create mocks.
		mrepo := &mservice.FailureRepository{}
		mrepo.On("GetAllByNode", mock.Anything).Once().Return(test.failures)

		// Create the service.
		fss := service.NewFailureStatus(mrepo, log.Dummy)

		// Get & check.
		gotFs := fss.GetNodeExpectedEnabledFailures("test")
		assert.Equal(test.expFailures, gotFs)
	}
}

func TestGetNodeExpectedDisabledFailures(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		failures    []*model.Failure
		expFailures []*model.Failure
	}{
		{
			failures: []*model.Failure{
				&model.Failure{ID: "f1", ExpectedState: types.EnabledFailureState},
			},
			expFailures: []*model.Failure{},
		},
		{
			failures: []*model.Failure{
				&model.Failure{ID: "f1", ExpectedState: types.DisabledFailureState},
			},
			expFailures: []*model.Failure{
				&model.Failure{ID: "f1", ExpectedState: types.DisabledFailureState},
			},
		},
		{
			failures: []*model.Failure{
				&model.Failure{ID: "f1", ExpectedState: types.DisabledFailureState},
				&model.Failure{ID: "f2", ExpectedState: types.EnabledFailureState},
				&model.Failure{ID: "f3", ExpectedState: types.RevertingFailureState},
				&model.Failure{ID: "f4", ExpectedState: types.EnabledFailureState},
			},
			expFailures: []*model.Failure{
				&model.Failure{ID: "f1", ExpectedState: types.DisabledFailureState},
			},
		},
	}

	for _, test := range tests {
		// Create mocks.
		mrepo := &mservice.FailureRepository{}
		mrepo.On("GetAllByNode", mock.Anything).Once().Return(test.failures)

		// Create the service.
		fss := service.NewFailureStatus(mrepo, log.Dummy)

		// Get & check.
		gotFs := fss.GetNodeExpectedDisabledFailures("test")
		assert.Equal(test.expFailures, gotFs)
	}
}

func TestGetFailure(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		expFailure *model.Failure
		expErr     bool
	}{
		{&model.Failure{ID: "test1"}, false},
		{&model.Failure{ID: "test2"}, true},
		{&model.Failure{ID: "test3"}, false},
	}

	for _, test := range tests {
		var err error

		// Create mocks.
		mrepo := &mservice.FailureRepository{}
		mrepo.On("Get", mock.Anything).Once().Return(test.expFailure, !test.expErr)

		// Create the service.
		fss := service.NewFailureStatus(mrepo, log.Dummy)

		// Get & check.
		f, err := fss.GetFailure(test.expFailure.ID)
		if test.expErr {
			assert.Error(err)
		} else if assert.NoError(err) {
			assert.Equal(test.expFailure, f)
		}
	}
}
