// Code generated by MockGen. DO NOT EDIT.
// Source: internal/common/pipeline/pipeline.go

// Package mock_pipeline is a generated GoMock package.
package mock_pipeline

import (
	handler "authstore/internal/common/http/handler"
	pipeline "authstore/internal/common/pipeline"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockHandle is a mock of Handle interface.
type MockHandle struct {
	ctrl     *gomock.Controller
	recorder *MockHandleMockRecorder
}

// MockHandleMockRecorder is the mock recorder for MockHandle.
type MockHandleMockRecorder struct {
	mock *MockHandle
}

// NewMockHandle creates a new mock instance.
func NewMockHandle(ctrl *gomock.Controller) *MockHandle {
	mock := &MockHandle{ctrl: ctrl}
	mock.recorder = &MockHandleMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHandle) EXPECT() *MockHandleMockRecorder {
	return m.recorder
}

// Handle mocks base method.
func (m *MockHandle) Handle(hctx *handler.HandleContext, next pipeline.Handle) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Handle", hctx, next)
	ret0, _ := ret[0].(error)
	return ret0
}

// Handle indicates an expected call of Handle.
func (mr *MockHandleMockRecorder) Handle(hctx, next interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Handle", reflect.TypeOf((*MockHandle)(nil).Handle), hctx, next)
}