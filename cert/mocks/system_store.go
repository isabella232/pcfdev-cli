// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/cert (interfaces: SystemStore)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of SystemStore interface
type MockSystemStore struct {
	ctrl     *gomock.Controller
	recorder *_MockSystemStoreRecorder
}

// Recorder for MockSystemStore (not exported)
type _MockSystemStoreRecorder struct {
	mock *MockSystemStore
}

func NewMockSystemStore(ctrl *gomock.Controller) *MockSystemStore {
	mock := &MockSystemStore{ctrl: ctrl}
	mock.recorder = &_MockSystemStoreRecorder{mock}
	return mock
}

func (_m *MockSystemStore) EXPECT() *_MockSystemStoreRecorder {
	return _m.recorder
}

func (_m *MockSystemStore) Store(_param0 string) error {
	ret := _m.ctrl.Call(_m, "Store", _param0)
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockSystemStoreRecorder) Store(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Store", arg0)
}

func (_m *MockSystemStore) Unstore() error {
	ret := _m.ctrl.Call(_m, "Unstore")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockSystemStoreRecorder) Unstore() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Unstore")
}
