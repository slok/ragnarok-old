// Code generated by mockery v1.0.0
package failurestatus

import context "golang.org/x/net/context"
import github_com_slok_ragnarok_grpc_failurestatus "github.com/slok/ragnarok/grpc/failurestatus"
import metadata "google.golang.org/grpc/metadata"
import mock "github.com/stretchr/testify/mock"

// FailureStatus_FailureStateListServer is an autogenerated mock type for the FailureStatus_FailureStateListServer type
type FailureStatus_FailureStateListServer struct {
	mock.Mock
}

// Context provides a mock function with given fields:
func (_m *FailureStatus_FailureStateListServer) Context() context.Context {
	ret := _m.Called()

	var r0 context.Context
	if rf, ok := ret.Get(0).(func() context.Context); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(context.Context)
		}
	}

	return r0
}

// RecvMsg provides a mock function with given fields: m
func (_m *FailureStatus_FailureStateListServer) RecvMsg(m interface{}) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Send provides a mock function with given fields: _a0
func (_m *FailureStatus_FailureStateListServer) Send(_a0 *github_com_slok_ragnarok_grpc_failurestatus.FailuresState) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*github_com_slok_ragnarok_grpc_failurestatus.FailuresState) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendHeader provides a mock function with given fields: _a0
func (_m *FailureStatus_FailureStateListServer) SendHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendMsg provides a mock function with given fields: m
func (_m *FailureStatus_FailureStateListServer) SendMsg(m interface{}) error {
	ret := _m.Called(m)

	var r0 error
	if rf, ok := ret.Get(0).(func(interface{}) error); ok {
		r0 = rf(m)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetHeader provides a mock function with given fields: _a0
func (_m *FailureStatus_FailureStateListServer) SetHeader(_a0 metadata.MD) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(metadata.MD) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SetTrailer provides a mock function with given fields: _a0
func (_m *FailureStatus_FailureStateListServer) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}
