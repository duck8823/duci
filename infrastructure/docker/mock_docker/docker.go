// Code generated by MockGen. DO NOT EDIT.
// Source: infrastructure/docker/docker.go

// Package mock_docker is a generated GoMock package.
package mock_docker

import (
	context "github.com/duck8823/duci/infrastructure/context"
	docker "github.com/duck8823/duci/infrastructure/docker"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockClient is a mock of Client interface
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Build mocks base method
func (m *MockClient) Build(ctx context.Context, file io.Reader, tag, dockerfile string) (docker.Logger, error) {
	ret := m.ctrl.Call(m, "Build", ctx, file, tag, dockerfile)
	ret0, _ := ret[0].(docker.Logger)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Build indicates an expected call of Build
func (mr *MockClientMockRecorder) Build(ctx, file, tag, dockerfile interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockClient)(nil).Build), ctx, file, tag, dockerfile)
}

// Run mocks base method
func (m *MockClient) Run(ctx context.Context, opts docker.RuntimeOptions, tag string, cmd ...string) (string, docker.Logger, error) {
	varargs := []interface{}{ctx, opts, tag}
	for _, a := range cmd {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Run", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(docker.Logger)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Run indicates an expected call of Run
func (mr *MockClientMockRecorder) Run(ctx, opts, tag interface{}, cmd ...interface{}) *gomock.Call {
	varargs := append([]interface{}{ctx, opts, tag}, cmd...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockClient)(nil).Run), varargs...)
}

// Rm mocks base method
func (m *MockClient) Rm(ctx context.Context, containerId string) error {
	ret := m.ctrl.Call(m, "Rm", ctx, containerId)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rm indicates an expected call of Rm
func (mr *MockClientMockRecorder) Rm(ctx, containerId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rm", reflect.TypeOf((*MockClient)(nil).Rm), ctx, containerId)
}

// Rmi mocks base method
func (m *MockClient) Rmi(ctx context.Context, tag string) error {
	ret := m.ctrl.Call(m, "Rmi", ctx, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rmi indicates an expected call of Rmi
func (mr *MockClientMockRecorder) Rmi(ctx, tag interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rmi", reflect.TypeOf((*MockClient)(nil).Rmi), ctx, tag)
}

// ExitCode mocks base method
func (m *MockClient) ExitCode(ctx context.Context, containerId string) (int64, error) {
	ret := m.ctrl.Call(m, "ExitCode", ctx, containerId)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExitCode indicates an expected call of ExitCode
func (mr *MockClientMockRecorder) ExitCode(ctx, containerId interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExitCode", reflect.TypeOf((*MockClient)(nil).ExitCode), ctx, containerId)
}
