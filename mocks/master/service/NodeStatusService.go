// Code generated by mockery v1.0.0
package service

import mock "github.com/stretchr/testify/mock"

import v1 "github.com/slok/ragnarok/api/cluster/v1"

// NodeStatusService is an autogenerated mock type for the NodeStatusService type
type NodeStatusService struct {
	mock.Mock
}

// Heartbeat provides a mock function with given fields: id, state
func (_m *NodeStatusService) Heartbeat(id string, state v1.NodeState) error {
	ret := _m.Called(id, state)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, v1.NodeState) error); ok {
		r0 = rf(id, state)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Register provides a mock function with given fields: id, labels
func (_m *NodeStatusService) Register(id string, labels map[string]string) error {
	ret := _m.Called(id, labels)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, map[string]string) error); ok {
		r0 = rf(id, labels)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
