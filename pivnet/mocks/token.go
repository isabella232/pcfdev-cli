// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/pivnet (interfaces: PivnetToken)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of PivnetToken interface
type MockPivnetToken struct {
	ctrl     *gomock.Controller
	recorder *_MockPivnetTokenRecorder
}

// Recorder for MockPivnetToken (not exported)
type _MockPivnetTokenRecorder struct {
	mock *MockPivnetToken
}

func NewMockPivnetToken(ctrl *gomock.Controller) *MockPivnetToken {
	mock := &MockPivnetToken{ctrl: ctrl}
	mock.recorder = &_MockPivnetTokenRecorder{mock}
	return mock
}

func (_m *MockPivnetToken) EXPECT() *_MockPivnetTokenRecorder {
	return _m.recorder
}

func (_m *MockPivnetToken) Destroy() error {
	ret := _m.ctrl.Call(_m, "Destroy")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockPivnetTokenRecorder) Destroy() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Destroy")
}

func (_m *MockPivnetToken) Get() (string, error) {
	ret := _m.ctrl.Call(_m, "Get")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockPivnetTokenRecorder) Get() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Get")
}
