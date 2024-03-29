// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/metadata"

	v1 "github.com/elizabeth-dev/ACME_Test/pkg/api/v1"
)

// UserService_GetUsersServer is an autogenerated mock type for the UserService_GetUsersServer type
type UserService_GetUsersServer struct {
	mock.Mock
}

// Context provides a mock function with given fields:
func (_m *UserService_GetUsersServer) Context() context.Context {
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
func (_m *UserService_GetUsersServer) RecvMsg(m interface{}) error {
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
func (_m *UserService_GetUsersServer) Send(_a0 *v1.User) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(*v1.User) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// SendHeader provides a mock function with given fields: _a0
func (_m *UserService_GetUsersServer) SendHeader(_a0 metadata.MD) error {
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
func (_m *UserService_GetUsersServer) SendMsg(m interface{}) error {
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
func (_m *UserService_GetUsersServer) SetHeader(_a0 metadata.MD) error {
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
func (_m *UserService_GetUsersServer) SetTrailer(_a0 metadata.MD) {
	_m.Called(_a0)
}

type mockConstructorTestingTNewUserService_GetUsersServer interface {
	mock.TestingT
	Cleanup(func())
}

// NewUserService_GetUsersServer creates a new instance of UserService_GetUsersServer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUserService_GetUsersServer(t mockConstructorTestingTNewUserService_GetUsersServer) *UserService_GetUsersServer {
	mock := &UserService_GetUsersServer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
