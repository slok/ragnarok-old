// Code generated by mockery v1.0.0
package grpc

import context "golang.org/x/net/context"
import grpc "google.golang.org/grpc"
import mock "github.com/stretchr/testify/mock"
import nodestatus "github.com/slok/ragnarok/grpc/nodestatus"

// NodeStatusClient is an autogenerated mock type for the NodeStatusClient type
type NodeStatusClient struct {
	mock.Mock
}

// Register provides a mock function with given fields: ctx, in, opts
func (_m *NodeStatusClient) Register(ctx context.Context, in *nodestatus.Node, opts ...grpc.CallOption) (*nodestatus.RegisteredResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *nodestatus.RegisteredResponse
	if rf, ok := ret.Get(0).(func(context.Context, *nodestatus.Node, ...grpc.CallOption) *nodestatus.RegisteredResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*nodestatus.RegisteredResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *nodestatus.Node, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
