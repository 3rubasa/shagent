// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/3rubasa/shagent/pkg/drivers/relay/asyncdecorator (interfaces: DeviceAPI)

// Package mockdeviceapi is a generated GoMock package.
package mockdeviceapi

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockDeviceAPI is a mock of DeviceAPI interface.
type MockDeviceAPI struct {
	ctrl     *gomock.Controller
	recorder *MockDeviceAPIMockRecorder
}

// MockDeviceAPIMockRecorder is the mock recorder for MockDeviceAPI.
type MockDeviceAPIMockRecorder struct {
	mock *MockDeviceAPI
}

// NewMockDeviceAPI creates a new mock instance.
func NewMockDeviceAPI(ctrl *gomock.Controller) *MockDeviceAPI {
	mock := &MockDeviceAPI{ctrl: ctrl}
	mock.recorder = &MockDeviceAPIMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDeviceAPI) EXPECT() *MockDeviceAPIMockRecorder {
	return m.recorder
}

// GetState mocks base method.
func (m *MockDeviceAPI) GetState() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetState")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetState indicates an expected call of GetState.
func (mr *MockDeviceAPIMockRecorder) GetState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetState", reflect.TypeOf((*MockDeviceAPI)(nil).GetState))
}

// TurnOff mocks base method.
func (m *MockDeviceAPI) TurnOff() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TurnOff")
	ret0, _ := ret[0].(error)
	return ret0
}

// TurnOff indicates an expected call of TurnOff.
func (mr *MockDeviceAPIMockRecorder) TurnOff() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TurnOff", reflect.TypeOf((*MockDeviceAPI)(nil).TurnOff))
}

// TurnOn mocks base method.
func (m *MockDeviceAPI) TurnOn() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TurnOn")
	ret0, _ := ret[0].(error)
	return ret0
}

// TurnOn indicates an expected call of TurnOn.
func (mr *MockDeviceAPIMockRecorder) TurnOn() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TurnOn", reflect.TypeOf((*MockDeviceAPI)(nil).TurnOn))
}
