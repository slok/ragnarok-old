// Code generated by mockery v1.0.0
package client

import client "github.com/slok/ragnarok/node/client"
import failure "github.com/slok/ragnarok/chaos/failure"
import mock "github.com/stretchr/testify/mock"

// Failure is an autogenerated mock type for the Failure type
type Failure struct {
	mock.Mock
}

// GetFailure provides a mock function with given fields: id
func (_m *Failure) GetFailure(id string) (*failure.Failure, error) {
	ret := _m.Called(id)

	var r0 *failure.Failure
	if rf, ok := ret.Get(0).(func(string) *failure.Failure); ok {
		r0 = rf(id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*failure.Failure)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ProcessFailureStateStreaming provides a mock function with given fields: nodeID, handler, stopCh
func (_m *Failure) ProcessFailureStateStreaming(nodeID string, handler client.FailureStateHandler, stopCh <-chan struct{}) error {
	ret := _m.Called(nodeID, handler, stopCh)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, client.FailureStateHandler, <-chan struct{}) error); ok {
		r0 = rf(nodeID, handler, stopCh)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
