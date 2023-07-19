// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	book "github.com/vladimirpekarski/wordofwisdom/internal/book"
)

// Booker is an autogenerated mock type for the Booker type
type Booker struct {
	mock.Mock
}

// RandomQuote provides a mock function with given fields:
func (_m *Booker) RandomQuote() book.Record {
	ret := _m.Called()

	var r0 book.Record
	if rf, ok := ret.Get(0).(func() book.Record); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(book.Record)
	}

	return r0
}

type mockConstructorTestingTNewBooker interface {
	mock.TestingT
	Cleanup(func())
}

// NewBooker creates a new instance of Booker. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewBooker(t mockConstructorTestingTNewBooker) *Booker {
	mock := &Booker{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
