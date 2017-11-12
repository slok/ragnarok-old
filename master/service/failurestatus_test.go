package service_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service"
	mrepository "github.com/slok/ragnarok/mocks/master/service/repository"
)

type testNodeFailures map[string][]*v1.Failure

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
					&v1.Failure{Metadata: v1.FailureMetadata{ID: "f11", NodeID: "node1"}},
				},
				"node2": {
					&v1.Failure{Metadata: v1.FailureMetadata{ID: "f21", NodeID: "node2"}},
					&v1.Failure{Metadata: v1.FailureMetadata{ID: "f21", NodeID: "node2"}},
				},
				"node3": {},
			},
		},
	}

	for _, test := range tests {
		// Create mocks.
		mrepo := &mrepository.Failure{}

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
		failures    []*v1.Failure
		expFailures []*v1.Failure
	}{
		{
			failures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
				},
			},
			expFailures: []*v1.Failure{},
		},
		{
			failures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
			},
			expFailures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
			},
		},
		{
			failures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
				},
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f2"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f3"},
					Status:   v1.FailureStatus{ExpectedState: v1.RevertingFailureState},
				},
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f4"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
			},
			expFailures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f2"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f4"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
			},
		},
	}

	for _, test := range tests {
		// Create mocks.
		mrepo := &mrepository.Failure{}
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
		failures    []*v1.Failure
		expFailures []*v1.Failure
	}{
		{
			failures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
			},
			expFailures: []*v1.Failure{},
		},
		{
			failures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
				},
			},
			expFailures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
				},
			},
		},
		{
			failures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
				},
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f2"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f3"},
					Status:   v1.FailureStatus{ExpectedState: v1.RevertingFailureState},
				},
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f4"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
			},
			expFailures: []*v1.Failure{
				&v1.Failure{
					Metadata: v1.FailureMetadata{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
				},
			},
		},
	}

	for _, test := range tests {
		// Create mocks.
		mrepo := &mrepository.Failure{}
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
		expFailure *v1.Failure
		expErr     bool
	}{
		{
			expFailure: &v1.Failure{Metadata: v1.FailureMetadata{ID: "test1"}},
			expErr:     false,
		},
		{
			expFailure: &v1.Failure{Metadata: v1.FailureMetadata{ID: "test2"}},
			expErr:     true,
		},
		{
			expFailure: &v1.Failure{Metadata: v1.FailureMetadata{ID: "test3"}},
			expErr:     false,
		},
	}

	for _, test := range tests {
		var err error

		// Create mocks.
		mrepo := &mrepository.Failure{}
		mrepo.On("Get", mock.Anything).Once().Return(test.expFailure, !test.expErr)

		// Create the service.
		fss := service.NewFailureStatus(mrepo, log.Dummy)

		// Get & check.
		f, err := fss.GetFailure(test.expFailure.Metadata.ID)
		if test.expErr {
			assert.Error(err)
		} else if assert.NoError(err) {
			assert.Equal(test.expFailure, f)
		}
	}
}
