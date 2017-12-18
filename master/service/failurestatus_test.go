package service_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/slok/ragnarok/api"
	"github.com/slok/ragnarok/api/chaos/v1"
	"github.com/slok/ragnarok/log"
	"github.com/slok/ragnarok/master/service"
	mclichaosv1 "github.com/slok/ragnarok/mocks/client/api/chaos/v1"
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
					&v1.Failure{Metadata: api.ObjectMeta{ID: "f11"}},
				},
				"node2": {
					&v1.Failure{Metadata: api.ObjectMeta{ID: "f21"}},
					&v1.Failure{Metadata: api.ObjectMeta{ID: "f21"}},
				},
				"node3": {},
			},
		},
	}

	for _, test := range tests {
		// Create mocks.
		mcli := &mclichaosv1.FailureClientInterface{}

		// Mock the call.
		for nID, expfs := range test.expectedFailures {
			opts := api.ListOptions{
				LabelSelector: map[string]string{
					api.LabelNode: nID,
				},
			}
			fsList := v1.NewFailureList(expfs, "")
			mcli.On("List", opts).Once().Return(&fsList, nil)
		}

		// Create the service.
		fs := service.NewFailureStatus(mcli, log.Dummy)

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
		failures    *v1.FailureList
		expFailures []*v1.Failure
	}{
		{
			failures: &v1.FailureList{
				Items: []*v1.Failure{
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f1"},
						Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
					},
				},
			},
			expFailures: []*v1.Failure{},
		},
		{
			failures: &v1.FailureList{
				Items: []*v1.Failure{
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f1"},
						Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
					},
				},
			},
			expFailures: []*v1.Failure{
				&v1.Failure{
					Metadata: api.ObjectMeta{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
			},
		},
		{
			failures: &v1.FailureList{
				Items: []*v1.Failure{
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f1"},
						Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
					},
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f2"},
						Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
					},
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f3"},
						Status:   v1.FailureStatus{ExpectedState: v1.RevertingFailureState},
					},
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f4"},
						Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
					},
				},
			},
			expFailures: []*v1.Failure{
				&v1.Failure{
					Metadata: api.ObjectMeta{ID: "f2"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
				&v1.Failure{
					Metadata: api.ObjectMeta{ID: "f4"},
					Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
				},
			},
		},
	}

	for _, test := range tests {
		// Create mocks.
		mcli := &mclichaosv1.FailureClientInterface{}
		mcli.On("List", mock.Anything).Once().Return(test.failures, nil)

		// Create the service.
		fss := service.NewFailureStatus(mcli, log.Dummy)

		// Get & check.
		gotFs := fss.GetNodeExpectedEnabledFailures("test")
		assert.Equal(test.expFailures, gotFs)
	}
}

func TestGetNodeExpectedDisabledFailures(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		failures    *v1.FailureList
		expFailures []*v1.Failure
	}{
		{
			failures: &v1.FailureList{
				Items: []*v1.Failure{
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f1"},
						Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
					},
				},
			},
			expFailures: []*v1.Failure{},
		},
		{
			failures: &v1.FailureList{
				Items: []*v1.Failure{
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f1"},
						Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
					},
				},
			},
			expFailures: []*v1.Failure{
				&v1.Failure{
					Metadata: api.ObjectMeta{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
				},
			},
		},
		{
			failures: &v1.FailureList{
				Items: []*v1.Failure{
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f1"},
						Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
					},
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f2"},
						Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
					},
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f3"},
						Status:   v1.FailureStatus{ExpectedState: v1.RevertingFailureState},
					},
					&v1.Failure{
						Metadata: api.ObjectMeta{ID: "f4"},
						Status:   v1.FailureStatus{ExpectedState: v1.EnabledFailureState},
					},
				},
			},
			expFailures: []*v1.Failure{
				&v1.Failure{
					Metadata: api.ObjectMeta{ID: "f1"},
					Status:   v1.FailureStatus{ExpectedState: v1.DisabledFailureState},
				},
			},
		},
	}

	for _, test := range tests {
		// Create mocks.
		mcli := &mclichaosv1.FailureClientInterface{}
		mcli.On("List", mock.Anything).Once().Return(test.failures, nil)

		// Create the service.
		fss := service.NewFailureStatus(mcli, log.Dummy)

		// Get & check.
		gotFs := fss.GetNodeExpectedDisabledFailures("test")
		assert.Equal(test.expFailures, gotFs)
	}
}

func TestGetFailure(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		expFailure *v1.Failure
		getError   bool
		expErr     bool
	}{
		{
			expFailure: &v1.Failure{Metadata: api.ObjectMeta{ID: "test1"}},
			getError:   false,
			expErr:     false,
		},
		{
			expFailure: &v1.Failure{Metadata: api.ObjectMeta{ID: "test2"}},
			getError:   true,
			expErr:     true,
		},
		{
			expFailure: &v1.Failure{Metadata: api.ObjectMeta{ID: "test3"}},
			getError:   false,
			expErr:     false,
		},
	}

	for _, test := range tests {
		var getError error
		if test.getError {
			getError = errors.New("wanted error")
		}
		// Create mocks.
		mcli := &mclichaosv1.FailureClientInterface{}
		mcli.On("Get", mock.Anything).Once().Return(test.expFailure, getError)

		// Create the service.
		fss := service.NewFailureStatus(mcli, log.Dummy)

		// Get & check.
		f, err := fss.GetFailure(test.expFailure.Metadata.ID)
		if test.expErr {
			assert.Error(err)
		} else if assert.NoError(err) {
			assert.Equal(test.expFailure, f)
		}
	}
}
