// Code generated by mockery v1.0.0
package service

import mock "github.com/stretchr/testify/mock"

import time "time"

// Status is an autogenerated mock type for the Status type
type Status struct {
	mock.Mock
}

// DeregisterOnMaster provides a mock function with given fields:
func (_m *Status) DeregisterOnMaster() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// RegisterOnMaster provides a mock function with given fields:
func (_m *Status) RegisterOnMaster() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// StartHeartbeat provides a mock function with given fields: interval
func (_m *Status) StartHeartbeat(interval time.Duration) (chan error, error) {
	ret := _m.Called(interval)

	var r0 chan error
	if rf, ok := ret.Get(0).(func(time.Duration) chan error); ok {
		r0 = rf(interval)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(chan error)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(time.Duration) error); ok {
		r1 = rf(interval)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// State provides a mock function with given fields:
func (_m *Status) State() bool {
	ret := _m.Called()

	var r0 bool
	if rf, ok := ret.Get(0).(func() bool); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

// StopHeartbeat provides a mock function with given fields:
func (_m *Status) StopHeartbeat() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
