// Code generated by mockery v1.0.0
package types

import failure "github.com/slok/ragnarok/grpc/failurestatus"
import mock "github.com/stretchr/testify/mock"

import v1 "github.com/slok/ragnarok/api/chaos/v1"

// FailureParser is an autogenerated mock type for the FailureParser type
type FailureParser struct {
	mock.Mock
}

// FailureToPB provides a mock function with given fields: fl
func (_m *FailureParser) FailureToPB(fl *v1.Failure) (*failure.Failure, error) {
	ret := _m.Called(fl)

	var r0 *failure.Failure
	if rf, ok := ret.Get(0).(func(*v1.Failure) *failure.Failure); ok {
		r0 = rf(fl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*failure.Failure)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*v1.Failure) error); ok {
		r1 = rf(fl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PBToFailure provides a mock function with given fields: fl
func (_m *FailureParser) PBToFailure(fl *failure.Failure) (*v1.Failure, error) {
	ret := _m.Called(fl)

	var r0 *v1.Failure
	if rf, ok := ret.Get(0).(func(*failure.Failure) *v1.Failure); ok {
		r0 = rf(fl)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*v1.Failure)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*failure.Failure) error); ok {
		r1 = rf(fl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}