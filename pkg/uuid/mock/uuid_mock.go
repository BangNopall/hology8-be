package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockUUIDInterface is a mock of UUIDInterface interface.
type MockUUIDInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUUIDInterfaceMockRecorder
}

// MockUUIDInterfaceMockRecorder is the mock recorder for MockUUIDInterface.
type MockUUIDInterfaceMockRecorder struct {
	mock *MockUUIDInterface
}

// NewMockUUIDInterface creates a new mock instance.
func NewMockUUIDInterface(ctrl *gomock.Controller) *MockUUIDInterface {
	mock := &MockUUIDInterface{ctrl: ctrl}
	mock.recorder = &MockUUIDInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUUIDInterface) EXPECT() *MockUUIDInterfaceMockRecorder {
	return m.recorder
}

// New mocks base method.
func (m *MockUUIDInterface) New() (uuid.UUID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "New")
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// New indicates an expected call of New.
func (mr *MockUUIDInterfaceMockRecorder) New() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "New", reflect.TypeOf((*MockUUIDInterface)(nil).New))
}
