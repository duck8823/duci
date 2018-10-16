// Code generated by MockGen. DO NOT EDIT.
// Source: application/service/docker/docker.go

// Package mock_docker is a generated GoMock package.
package mock_docker

import (
	context "context"
	docker "github.com/duck8823/duci/application/service/docker"
	gomock "github.com/golang/mock/gomock"
	io "io"
	reflect "reflect"
)

// MockService is a mock of Service interface
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// Build mocks base method
func (m *MockService) Build(ctx context.Context, file io.Reader, tag docker.Tag, dockerfile docker.Dockerfile) (docker.Log, error) {
	ret := m.ctrl.Call(m, "Build", ctx, file, tag, dockerfile)
	ret0, _ := ret[0].(docker.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Build indicates an expected call of Build
func (mr *MockServiceMockRecorder) Build(ctx, file, tag, dockerfile interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Build", reflect.TypeOf((*MockService)(nil).Build), ctx, file, tag, dockerfile)
}

// Run mocks base method
func (m *MockService) Run(ctx context.Context, opts docker.RuntimeOptions, tag docker.Tag, cmd docker.Command) (docker.ContainerID, docker.Log, error) {
	ret := m.ctrl.Call(m, "Run", ctx, opts, tag, cmd)
	ret0, _ := ret[0].(docker.ContainerID)
	ret1, _ := ret[1].(docker.Log)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Run indicates an expected call of Run
func (mr *MockServiceMockRecorder) Run(ctx, opts, tag, cmd interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Run", reflect.TypeOf((*MockService)(nil).Run), ctx, opts, tag, cmd)
}

// Rm mocks base method
func (m *MockService) Rm(ctx context.Context, containerID docker.ContainerID) error {
	ret := m.ctrl.Call(m, "Rm", ctx, containerID)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rm indicates an expected call of Rm
func (mr *MockServiceMockRecorder) Rm(ctx, containerID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rm", reflect.TypeOf((*MockService)(nil).Rm), ctx, containerID)
}

// Rmi mocks base method
func (m *MockService) Rmi(ctx context.Context, tag docker.Tag) error {
	ret := m.ctrl.Call(m, "Rmi", ctx, tag)
	ret0, _ := ret[0].(error)
	return ret0
}

// Rmi indicates an expected call of Rmi
func (mr *MockServiceMockRecorder) Rmi(ctx, tag interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rmi", reflect.TypeOf((*MockService)(nil).Rmi), ctx, tag)
}

// ExitCode mocks base method
func (m *MockService) ExitCode(ctx context.Context, containerID docker.ContainerID) (docker.ExitCode, error) {
	ret := m.ctrl.Call(m, "ExitCode", ctx, containerID)
	ret0, _ := ret[0].(docker.ExitCode)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExitCode indicates an expected call of ExitCode
func (mr *MockServiceMockRecorder) ExitCode(ctx, containerID interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExitCode", reflect.TypeOf((*MockService)(nil).ExitCode), ctx, containerID)
}

// Status mocks base method
func (m *MockService) Status() error {
	ret := m.ctrl.Call(m, "Status")
	ret0, _ := ret[0].(error)
	return ret0
}

// Status indicates an expected call of Status
func (mr *MockServiceMockRecorder) Status() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Status", reflect.TypeOf((*MockService)(nil).Status))
}
