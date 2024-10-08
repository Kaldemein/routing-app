// Code generated by MockGen. DO NOT EDIT.
// Source: types.go

// Package mock_api is a generated GoMock package.
package mock_api

import (
        reflect "reflect"
        gomock "github.com/golang/mock/gomock"
)

// MockMessageSender is a mock of MessageSender interface.
type MockMessageSender struct {
        ctrl     *gomock.Controller
        recorder *MockMessageSenderMockRecorder
}

// MockMessageSenderMockRecorder is the mock recorder for MockMessageSender.
type MockMessageSenderMockRecorder struct {
        mock *MockMessageSender
}

// NewMockMessageSender creates a new mock instance.
func NewMockMessageSender(ctrl *gomock.Controller) *MockMessageSender {
        mock := &MockMessageSender{ctrl: ctrl}
        mock.recorder = &MockMessageSenderMockRecorder{mock}
        return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMessageSender) EXPECT() *MockMessageSenderMockRecorder {
        return m.recorder
}

// SendVerifyMessageToRMQ mocks base method.
func (m *MockMessageSender) SendVerifyMessageToRMQ(arg0 string) (string, error) {
        m.ctrl.T.Helper()
        ret := m.ctrl.Call(m, "SendVerifyMessageToRMQ", arg0)
        ret0, _ := ret[0].(string)
        ret1, _ := ret[1].(error)
        return ret0, ret1
}

// SendVerifyMessageToRMQ indicates an expected call of SendVerifyMessageToRMQ.
func (mr *MockMessageSenderMockRecorder) SendVerifyMessageToRMQ(arg0 interface{}) *gomock.Call {
        mr.mock.ctrl.T.Helper()
        return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendVerifyMessageToRMQ", reflect.TypeOf((*MockMessageSender)(nil).SendVerifyMessageToRMQ), arg0)
}