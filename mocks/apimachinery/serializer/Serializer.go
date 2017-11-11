// Code generated by mockery v1.0.0
package serializer

import api "github.com/slok/ragnarok/api"
import mock "github.com/stretchr/testify/mock"

// Serializer is an autogenerated mock type for the Serializer type
type Serializer struct {
	mock.Mock
}

// Decode provides a mock function with given fields: data
func (_m *Serializer) Decode(data interface{}) (api.Object, error) {
	ret := _m.Called(data)

	var r0 api.Object
	if rf, ok := ret.Get(0).(func(interface{}) api.Object); ok {
		r0 = rf(data)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(api.Object)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(interface{}) error); ok {
		r1 = rf(data)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Encode provides a mock function with given fields: obj, out
func (_m *Serializer) Encode(obj api.Object, out interface{}) error {
	ret := _m.Called(obj, out)

	var r0 error
	if rf, ok := ret.Get(0).(func(api.Object, interface{}) error); ok {
		r0 = rf(obj, out)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}