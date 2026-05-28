package mock

import (
	reflect "reflect"
	time "time"

	gomock "github.com/golang/mock/gomock"
)

// MockTimeInterface is a mock of TimeInterface interface.
type MockTimeInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTimeInterfaceMockRecorder
}

// MockTimeInterfaceMockRecorder is the mock recorder for MockTimeInterface.
type MockTimeInterfaceMockRecorder struct {
	mock *MockTimeInterface
}

// NewMockTimeInterface creates a new mock instance.
func NewMockTimeInterface(ctrl *gomock.Controller) *MockTimeInterface {
	mock := &MockTimeInterface{ctrl: ctrl}
	mock.recorder = &MockTimeInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTimeInterface) EXPECT() *MockTimeInterfaceMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockTimeInterface) Add(duration time.Duration) time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", duration)
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockTimeInterfaceMockRecorder) Add(duration interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockTimeInterface)(nil).Add), duration)
}

// Now mocks base method.
func (m *MockTimeInterface) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockTimeInterfaceMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockTimeInterface)(nil).Now))
}
