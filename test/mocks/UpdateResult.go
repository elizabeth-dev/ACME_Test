// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import "github.com/stretchr/testify/mock"

// UpdateResult is an autogenerated mock type for the UpdateResult type
type UpdateResult struct {
	mock.Mock
}

type mockConstructorTestingTNewUpdateResult interface {
	mock.TestingT
	Cleanup(func())
}

// NewUpdateResult creates a new instance of UpdateResult. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewUpdateResult(t mockConstructorTestingTNewUpdateResult) *UpdateResult {
	mock := &UpdateResult{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}