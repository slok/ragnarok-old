// Code generated by mockery v1.0.0
package client

import failure "github.com/slok/ragnarok/failure"
import mock "github.com/stretchr/testify/mock"

// FailureStateHandler is an autogenerated mock type for the FailureStateHandler type
type FailureStateHandler struct {
	mock.Mock
}

// ProcessFailureStates provides a mock function with given fields: failures
func (_m *FailureStateHandler) ProcessFailureStates(failures []*failure.Failure) error {
	ret := _m.Called(failures)

	var r0 error
	if rf, ok := ret.Get(0).(func([]*failure.Failure) error); ok {
		r0 = rf(failures)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}