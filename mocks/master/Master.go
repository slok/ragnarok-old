// Code generated by mockery v1.0.0
package master

import mock "github.com/stretchr/testify/mock"

// Master is an autogenerated mock type for the Master type
type Master struct {
	mock.Mock
}

// RegisterNode provides a mock function with given fields: id, tags
func (_m *Master) RegisterNode(id string, tags map[string]string) error {
	ret := _m.Called(id, tags)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, map[string]string) error); ok {
		r0 = rf(id, tags)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
