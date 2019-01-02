// Code generated by MockGen. DO NOT EDIT.
// Source: application/service/executor/executor.go

// Package mock_executor is a generated GoMock package.
package mock_executor

import (
	context "context"
	job "github.com/duck8823/duci/domain/model/job"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockExecutor is a mock of Executor interface
type MockExecutor struct {
	ctrl     *gomock.Controller
	recorder *MockExecutorMockRecorder
}

// MockExecutorMockRecorder is the mock recorder for MockExecutor
type MockExecutorMockRecorder struct {
	mock *MockExecutor
}

// NewMockExecutor creates a new mock instance
func NewMockExecutor(ctrl *gomock.Controller) *MockExecutor {
	mock := &MockExecutor{ctrl: ctrl}
	mock.recorder = &MockExecutorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockExecutor) EXPECT() *MockExecutorMockRecorder {
	return m.recorder
}

// Execute mocks base method
func (m *MockExecutor) Execute(ctx context.Context, target job.Target, cmd ...string) error {
	varargs := []interface{}{ctx, target}
	for _, a := range cmd {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Execute", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Execute indicates an expected call of Execute
func (mr *MockExecutorMockRecorder) Execute(ctx, target interface{}, cmd ...interface{}) *gomock.Call {
	varargs := append([]interface{}{ctx, target}, cmd...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockExecutor)(nil).Execute), varargs...)
}
