// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	book "github.com/vladimirpekarski/wordofwisdom/internal/book"
)

// RandomRecordGetter is an autogenerated mock type for the RandomRecordGetter type
type RandomRecordGetter struct {
	mock.Mock
}

// GetRandomRecord provides a mock function with given fields:
func (_m *RandomRecordGetter) GetRandomRecord() book.Record {
	ret := _m.Called()

	var r0 book.Record
	if rf, ok := ret.Get(0).(func() book.Record); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(book.Record)
	}

	return r0
}

type mockConstructorTestingTNewRandomRecordGetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewRandomRecordGetter creates a new instance of RandomRecordGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRandomRecordGetter(t mockConstructorTestingTNewRandomRecordGetter) *RandomRecordGetter {
	mock := &RandomRecordGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
