// Code generated by mockery v2.52.2. DO NOT EDIT.

package mocktasks

import (
	asynq "github.com/hibiken/asynq"
	mock "github.com/stretchr/testify/mock"
)

// Client is an autogenerated mock type for the Client type
type Client struct {
	mock.Mock
}

type Client_Expecter struct {
	mock *mock.Mock
}

func (_m *Client) EXPECT() *Client_Expecter {
	return &Client_Expecter{mock: &_m.Mock}
}

// Enqueue provides a mock function with given fields: task, opts
func (_m *Client) Enqueue(task *asynq.Task, opts ...asynq.Option) (*asynq.TaskInfo, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, task)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	if len(ret) == 0 {
		panic("no return value specified for Enqueue")
	}

	var r0 *asynq.TaskInfo
	var r1 error
	if rf, ok := ret.Get(0).(func(*asynq.Task, ...asynq.Option) (*asynq.TaskInfo, error)); ok {
		return rf(task, opts...)
	}
	if rf, ok := ret.Get(0).(func(*asynq.Task, ...asynq.Option) *asynq.TaskInfo); ok {
		r0 = rf(task, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*asynq.TaskInfo)
		}
	}

	if rf, ok := ret.Get(1).(func(*asynq.Task, ...asynq.Option) error); ok {
		r1 = rf(task, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Client_Enqueue_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Enqueue'
type Client_Enqueue_Call struct {
	*mock.Call
}

// Enqueue is a helper method to define mock.On call
//   - task *asynq.Task
//   - opts ...asynq.Option
func (_e *Client_Expecter) Enqueue(task interface{}, opts ...interface{}) *Client_Enqueue_Call {
	return &Client_Enqueue_Call{Call: _e.mock.On("Enqueue",
		append([]interface{}{task}, opts...)...)}
}

func (_c *Client_Enqueue_Call) Run(run func(task *asynq.Task, opts ...asynq.Option)) *Client_Enqueue_Call {
	_c.Call.Run(func(args mock.Arguments) {
		variadicArgs := make([]asynq.Option, len(args)-1)
		for i, a := range args[1:] {
			if a != nil {
				variadicArgs[i] = a.(asynq.Option)
			}
		}
		run(args[0].(*asynq.Task), variadicArgs...)
	})
	return _c
}

func (_c *Client_Enqueue_Call) Return(_a0 *asynq.TaskInfo, _a1 error) *Client_Enqueue_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *Client_Enqueue_Call) RunAndReturn(run func(*asynq.Task, ...asynq.Option) (*asynq.TaskInfo, error)) *Client_Enqueue_Call {
	_c.Call.Return(run)
	return _c
}

// NewClient creates a new instance of Client. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewClient(t interface {
	mock.TestingT
	Cleanup(func())
}) *Client {
	mock := &Client{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
