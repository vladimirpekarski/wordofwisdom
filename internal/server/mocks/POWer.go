// Code generated by mockery v2.28.2. DO NOT EDIT.

package mocks

import (
	mock "github.com/stretchr/testify/mock"
	message "github.com/vladimirpekarski/wordofwisdom/internal/message"
)

// POWer is an autogenerated mock type for the POWer type
type POWer struct {
	mock.Mock
}

// GenerateChallenge provides a mock function with given fields: difficulty
func (_m *POWer) GenerateChallenge(difficulty int) (message.Challenge, error) {
	ret := _m.Called(difficulty)

	var r0 message.Challenge
	var r1 error
	if rf, ok := ret.Get(0).(func(int) (message.Challenge, error)); ok {
		return rf(difficulty)
	}
	if rf, ok := ret.Get(0).(func(int) message.Challenge); ok {
		r0 = rf(difficulty)
	} else {
		r0 = ret.Get(0).(message.Challenge)
	}

	if rf, ok := ret.Get(1).(func(int) error); ok {
		r1 = rf(difficulty)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Validate provides a mock function with given fields: ch, sl
func (_m *POWer) Validate(ch message.Challenge, sl message.Solution) bool {
	ret := _m.Called(ch, sl)

	var r0 bool
	if rf, ok := ret.Get(0).(func(message.Challenge, message.Solution) bool); ok {
		r0 = rf(ch, sl)
	} else {
		r0 = ret.Get(0).(bool)
	}

	return r0
}

type mockConstructorTestingTNewPOWer interface {
	mock.TestingT
	Cleanup(func())
}

// NewPOWer creates a new instance of POWer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewPOWer(t mockConstructorTestingTNewPOWer) *POWer {
	mock := &POWer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}