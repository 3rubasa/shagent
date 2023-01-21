// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/3rubasa/shagent/controllers/watchdog (interfaces: InternetChecker)

// Package mockinternetchecker is a generated GoMock package.
package mockinternetchecker

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockInternetChecker is a mock of InternetChecker interface.
type MockInternetChecker struct {
	ctrl     *gomock.Controller
	recorder *MockInternetCheckerMockRecorder
}

// MockInternetCheckerMockRecorder is the mock recorder for MockInternetChecker.
type MockInternetCheckerMockRecorder struct {
	mock *MockInternetChecker
}

// NewMockInternetChecker creates a new mock instance.
func NewMockInternetChecker(ctrl *gomock.Controller) *MockInternetChecker {
	mock := &MockInternetChecker{ctrl: ctrl}
	mock.recorder = &MockInternetCheckerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockInternetChecker) EXPECT() *MockInternetCheckerMockRecorder {
	return m.recorder
}

// IsInternetAvailable mocks base method.
func (m *MockInternetChecker) IsInternetAvailable() (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsInternetAvailable")
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsInternetAvailable indicates an expected call of IsInternetAvailable.
func (mr *MockInternetCheckerMockRecorder) IsInternetAvailable() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsInternetAvailable", reflect.TypeOf((*MockInternetChecker)(nil).IsInternetAvailable))
}
