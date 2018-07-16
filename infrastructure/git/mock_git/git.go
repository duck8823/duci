// Code generated by MockGen. DO NOT EDIT.
// Source: infrastructure/git/git.go

// Package mock_git is a generated GoMock package.
package mock_git

import (
	context "github.com/duck8823/minimal-ci/infrastructure/context"
	gomock "github.com/golang/mock/gomock"
	plumbing "gopkg.in/src-d/go-git.v4/plumbing"
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

// Clone mocks base method
func (m *MockClient) Clone(ctx context.Context, dir, sshUrl, ref string) (plumbing.Hash, error) {
	ret := m.ctrl.Call(m, "Clone", ctx, dir, sshUrl, ref)
	ret0, _ := ret[0].(plumbing.Hash)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Clone indicates an expected call of Clone
func (mr *MockClientMockRecorder) Clone(ctx, dir, sshUrl, ref interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clone", reflect.TypeOf((*MockClient)(nil).Clone), ctx, dir, sshUrl, ref)
}
