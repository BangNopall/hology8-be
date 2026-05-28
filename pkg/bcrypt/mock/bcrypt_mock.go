package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockBcryptInterface is a mock of BcryptInterface interface.
type MockBcryptInterface struct {
	ctrl     *gomock.Controller
	recorder *MockBcryptInterfaceMockRecorder
}

// MockBcryptInterfaceMockRecorder is the mock recorder for MockBcryptInterface.
type MockBcryptInterfaceMockRecorder struct {
	mock *MockBcryptInterface
}

// NewMockBcryptInterface creates a new mock instance.
func NewMockBcryptInterface(ctrl *gomock.Controller) *MockBcryptInterface {
	mock := &MockBcryptInterface{ctrl: ctrl}
	mock.recorder = &MockBcryptInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockBcryptInterface) EXPECT() *MockBcryptInterfaceMockRecorder {
	return m.recorder
}

// Compare mocks base method.
func (m *MockBcryptInterface) Compare(password, hashed string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Compare", password, hashed)
	ret0, _ := ret[0].(bool)
	return ret0
}

// Compare indicates an expected call of Compare.
func (mr *MockBcryptInterfaceMockRecorder) Compare(password, hashed interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Compare", reflect.TypeOf((*MockBcryptInterface)(nil).Compare), password, hashed)
}

// Hash mocks base method.
func (m *MockBcryptInterface) Hash(plain string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Hash", plain)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Hash indicates an expected call of Hash.
func (mr *MockBcryptInterfaceMockRecorder) Hash(plain interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Hash", reflect.TypeOf((*MockBcryptInterface)(nil).Hash), plain)
}
