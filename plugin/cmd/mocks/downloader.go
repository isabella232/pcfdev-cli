// Automatically generated by MockGen. DO NOT EDIT!
// Source: github.com/pivotal-cf/pcfdev-cli/plugin/cmd (interfaces: Downloader)

package mocks

import (
	gomock "github.com/golang/mock/gomock"
)

// Mock of Downloader interface
type MockDownloader struct {
	ctrl     *gomock.Controller
	recorder *_MockDownloaderRecorder
}

// Recorder for MockDownloader (not exported)
type _MockDownloaderRecorder struct {
	mock *MockDownloader
}

func NewMockDownloader(ctrl *gomock.Controller) *MockDownloader {
	mock := &MockDownloader{ctrl: ctrl}
	mock.recorder = &_MockDownloaderRecorder{mock}
	return mock
}

func (_m *MockDownloader) EXPECT() *_MockDownloaderRecorder {
	return _m.recorder
}

func (_m *MockDownloader) Download() error {
	ret := _m.ctrl.Call(_m, "Download")
	ret0, _ := ret[0].(error)
	return ret0
}

func (_mr *_MockDownloaderRecorder) Download() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "Download")
}

func (_m *MockDownloader) IsOVACurrent() (bool, error) {
	ret := _m.ctrl.Call(_m, "IsOVACurrent")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (_mr *_MockDownloaderRecorder) IsOVACurrent() *gomock.Call {
	return _mr.mock.ctrl.RecordCall(_mr.mock, "IsOVACurrent")
}