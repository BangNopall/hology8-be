package mock

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockJwtInterface is a mock of JwtInterface interface.
type MockJwtInterface struct {
	ctrl     *gomock.Controller
	recorder *MockJwtInterfaceMockRecorder
}

// MockJwtInterfaceMockRecorder is the mock recorder for MockJwtInterface.
type MockJwtInterfaceMockRecorder struct {
	mock *MockJwtInterface
}

// NewMockJwtInterface creates a new mock instance.
func NewMockJwtInterface(ctrl *gomock.Controller) *MockJwtInterface {
	mock := &MockJwtInterface{ctrl: ctrl}
	mock.recorder = &MockJwtInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJwtInterface) EXPECT() *MockJwtInterfaceMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockJwtInterface) GenerateToken(userId uuid.UUID, entity, adminRole string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateToken", userId, entity, adminRole)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockJwtInterfaceMockRecorder) GenerateToken(userId, entity, adminRole interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockJwtInterface)(nil).GenerateToken), userId, entity, adminRole)
}

// ValidateToken mocks base method.
func (m *MockJwtInterface) ValidateToken(tokenString string) (uuid.UUID, string, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ValidateToken", tokenString)
	ret0, _ := ret[0].(uuid.UUID)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(string)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// ValidateToken indicates an expected call of ValidateToken.
func (mr *MockJwtInterfaceMockRecorder) ValidateToken(tokenString interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ValidateToken", reflect.TypeOf((*MockJwtInterface)(nil).ValidateToken), tokenString)
}
