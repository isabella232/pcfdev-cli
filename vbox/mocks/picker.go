// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/vbox (interfaces: NetworkPicker)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
	network "github.com/pivotal-cf/pcfdev-cli/network"
)

// Mock of NetworkPicker interface
type MockNetworkPicker struct {
	ctrl     *gomock.Controller
	recorder *_MockNetworkPickerRecorder
}

// Recorder for MockNetworkPicker (not exported)
type _MockNetworkPickerRecorder struct {
	mock *MockNetworkPicker
}

func NewMockNetworkPicker(ctrl *gomock.Controller) *MockNetworkPicker {
	mock := &MockNetworkPicker{ctrl: ctrl}
	mock.recorder = &_MockNetworkPickerRecorder{mock}
	return mock
}

func (_m *MockNetworkPicker) EXPECT() *_MockNetworkPickerRecorder {
	return _m.recorder
}

func (_m *MockNetworkPicker) SelectAvailableInterface(_param0 []*network.Interface) (*network.Interface, error) {
	ret := _m.ctrl.Call(_m, "SelectAvailableInterface", _param0)
	ret0, _ := ret[0].(*network.Interface)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockNetworkPickerRecorder) SelectAvailableInterface(arg0 interface{}) *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "SelectAvailableInterface", arg0)
}
