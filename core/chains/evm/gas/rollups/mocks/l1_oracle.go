// Code generated by mockery v2.35.4. DO NOT EDIT.

package mocks

import (
	context "context"

	assets "github.com/smartcontractkit/chainlink/v2/core/chains/evm/assets"

	mock "github.com/stretchr/testify/mock"
)

// L1Oracle is an autogenerated mock type for the L1Oracle type
type L1Oracle struct {
	mock.Mock
}

// Close provides a mock function with given fields:
func (_m *L1Oracle) Close() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GasPrice provides a mock function with given fields: ctx
func (_m *L1Oracle) GasPrice(ctx context.Context) (*assets.Wei, error) {
	ret := _m.Called(ctx)

	var r0 *assets.Wei
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*assets.Wei, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *assets.Wei); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*assets.Wei)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HealthReport provides a mock function with given fields:
func (_m *L1Oracle) HealthReport() map[string]error {
	ret := _m.Called()

	var r0 map[string]error
	if rf, ok := ret.Get(0).(func() map[string]error); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]error)
		}
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *L1Oracle) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Ready provides a mock function with given fields:
func (_m *L1Oracle) Ready() error {
	ret := _m.Called()

	var r0 error
	if rf, ok := ret.Get(0).(func() error); ok {
		r0 = rf()
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Start provides a mock function with given fields: _a0
func (_m *L1Oracle) Start(_a0 context.Context) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewL1Oracle creates a new instance of L1Oracle. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewL1Oracle(t interface {
	mock.TestingT
	Cleanup(func())
}) *L1Oracle {
	mock := &L1Oracle{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
