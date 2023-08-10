package baseapp

import (
	reflect "reflect"

	types "github.com/Finschia/finschia-sdk/types"
	gomock "github.com/golang/mock/gomock"
)

// MockXXXMessage is a mock of XXXMessage interface.
type MockXXXMessage struct {
	ctrl     *gomock.Controller
	recorder *MockXXXMessageMockRecorder
}

// MockXXXMessageMockRecorder is the mock recorder for MockXXXMessage.
type MockXXXMessageMockRecorder struct {
	mock *MockXXXMessage
}

// NewMockXXXMessage creates a new mock instance.
func NewMockXXXMessage(ctrl *gomock.Controller) *MockXXXMessage {
	mock := &MockXXXMessage{ctrl: ctrl}
	mock.recorder = &MockXXXMessageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockXXXMessage) EXPECT() *MockXXXMessageMockRecorder {
	return m.recorder
}

// GetSigners mocks base method.
func (m *MockXXXMessage) GetSigners() []types.AccAddress {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetSigners")
	ret0, _ := ret[0].([]types.AccAddress)
	return ret0
}

// GetSigners indicates an expected call of GetSigners.
func (mr *MockXXXMessageMockRecorder) GetSigners() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetSigners", reflect.TypeOf((*MockXXXMessage)(nil).GetSigners))
}

// ProtoMessage mocks base method.
func (m *MockXXXMessage) ProtoMessage() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "ProtoMessage")
}

// ProtoMessage indicates an expected call of ProtoMessage.
func (mr *MockXXXMessageMockRecorder) ProtoMessage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProtoMessage", reflect.TypeOf((*MockXXXMessage)(nil).ProtoMessage))
}

// Reset mocks base method.
func (m *MockXXXMessage) Reset() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Reset")
}

// Reset indicates an expected call of Reset.
func (mr *MockXXXMessageMockRecorder) Reset() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Reset", reflect.TypeOf((*MockXXXMessage)(nil).Reset))
}

// String mocks base method.
func (m *MockXXXMessage) String() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "String")
	ret0, _ := ret[0].(string)
	return ret0
}

// String indicates an expected call of String.
func (mr *MockXXXMessageMockRecorder) String() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "String", reflect.TypeOf((*MockXXXMessage)(nil).String))
}

// ValidateBasic mocks base method.
func (m *MockXXXMessage) ValidateBasic() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateBasic")
	ret0, _ := ret[0].(error)
	return ret0
}

// ValidateBasic indicates an expected call of ValidateBasic.
func (mr *MockXXXMessageMockRecorder) ValidateBasic() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateBasic", reflect.TypeOf((*MockXXXMessage)(nil).ValidateBasic))
}

// XXX_MessageName mocks base method.
func (m *MockXXXMessage) XXX_MessageName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "XXX_MessageName")
	ret0, _ := ret[0].(string)
	return ret0
}

// XXX_MessageName indicates an expected call of XXX_MessageName.
func (mr *MockXXXMessageMockRecorder) XXX_MessageName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "XXX_MessageName", reflect.TypeOf((*MockXXXMessage)(nil).XXX_MessageName))
}
