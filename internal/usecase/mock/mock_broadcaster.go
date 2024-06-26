// Code generated by mockery v2.43.0. DO NOT EDIT.

package ucmock

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockBroadcaster is an autogenerated mock type for the Broadcaster type
type MockBroadcaster struct {
	mock.Mock
}

type MockBroadcaster_Expecter struct {
	mock *mock.Mock
}

func (_m *MockBroadcaster) EXPECT() *MockBroadcaster_Expecter {
	return &MockBroadcaster_Expecter{mock: &_m.Mock}
}

// Broadcast provides a mock function with given fields: ctx, data
func (_m *MockBroadcaster) Broadcast(ctx context.Context, data []byte) error {
	ret := _m.Called(ctx, data)

	if len(ret) == 0 {
		panic("no return value specified for Broadcast")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, []byte) error); ok {
		r0 = rf(ctx, data)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// MockBroadcaster_Broadcast_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Broadcast'
type MockBroadcaster_Broadcast_Call struct {
	*mock.Call
}

// Broadcast is a helper method to define mock.On call
//   - ctx context.Context
//   - data []byte
func (_e *MockBroadcaster_Expecter) Broadcast(ctx interface{}, data interface{}) *MockBroadcaster_Broadcast_Call {
	return &MockBroadcaster_Broadcast_Call{Call: _e.mock.On("Broadcast", ctx, data)}
}

func (_c *MockBroadcaster_Broadcast_Call) Run(run func(ctx context.Context, data []byte)) *MockBroadcaster_Broadcast_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].([]byte))
	})
	return _c
}

func (_c *MockBroadcaster_Broadcast_Call) Return(_a0 error) *MockBroadcaster_Broadcast_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *MockBroadcaster_Broadcast_Call) RunAndReturn(run func(context.Context, []byte) error) *MockBroadcaster_Broadcast_Call {
	_c.Call.Return(run)
	return _c
}

// NewMockBroadcaster creates a new instance of MockBroadcaster. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewMockBroadcaster(t interface {
	mock.TestingT
	Cleanup(func())
}) *MockBroadcaster {
	mock := &MockBroadcaster{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
