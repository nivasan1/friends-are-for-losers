// Code generated by mockery v2.24.0. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// Filter is an autogenerated mock type for the Filter type
type Filter struct {
	mock.Mock
}

// Filter provides a mock function with given fields: ctx, address
func (_m *Filter) Filter(ctx context.Context, address string) (bool, error) {
	ret := _m.Called(ctx, address)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string) (bool, error)); ok {
		return rf(ctx, address)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string) bool); ok {
		r0 = rf(ctx, address)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, string) error); ok {
		r1 = rf(ctx, address)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewFilter interface {
	mock.TestingT
	Cleanup(func())
}

// NewFilter creates a new instance of Filter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFilter(t mockConstructorTestingTNewFilter) *Filter {
	mock := &Filter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
