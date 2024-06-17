// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/KartoonYoko/gophkeeper/internal/usecase/auth (interfaces: Storager)
//
// Generated by this command:
//
//	mockgen --destination=internal/usecase/auth/mocks/mock_auth_storage.go --package=mocks github.com/KartoonYoko/gophkeeper/internal/usecase/auth Storager
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"
	time "time"

	auth "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
	gomock "go.uber.org/mock/gomock"
)

// MockStorager is a mock of Storager interface.
type MockStorager struct {
	ctrl     *gomock.Controller
	recorder *MockStoragerMockRecorder
}

// MockStoragerMockRecorder is the mock recorder for MockStorager.
type MockStoragerMockRecorder struct {
	mock *MockStorager
}

// NewMockStorager creates a new mock instance.
func NewMockStorager(ctrl *gomock.Controller) *MockStorager {
	mock := &MockStorager{ctrl: ctrl}
	mock.recorder = &MockStoragerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorager) EXPECT() *MockStoragerMockRecorder {
	return m.recorder
}

// CreateRefreshToken mocks base method.
func (m *MockStorager) CreateRefreshToken(arg0 context.Context, arg1 *auth.CreateRefreshTokenRequestModel) (*auth.CreateRefreshTokenResponseModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRefreshToken", arg0, arg1)
	ret0, _ := ret[0].(*auth.CreateRefreshTokenResponseModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRefreshToken indicates an expected call of CreateRefreshToken.
func (mr *MockStoragerMockRecorder) CreateRefreshToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRefreshToken", reflect.TypeOf((*MockStorager)(nil).CreateRefreshToken), arg0, arg1)
}

// CreateUser mocks base method.
func (m *MockStorager) CreateUser(arg0 context.Context, arg1 *auth.CreateUserRequestModel) (*auth.CreateUserReqsponseModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", arg0, arg1)
	ret0, _ := ret[0].(*auth.CreateUserReqsponseModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockStoragerMockRecorder) CreateUser(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockStorager)(nil).CreateUser), arg0, arg1)
}

// GetRefreshToken mocks base method.
func (m *MockStorager) GetRefreshToken(arg0 context.Context, arg1 *auth.GetRefreshTokenRequestModel) (*auth.GetRefreshTokenResponseModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRefreshToken", arg0, arg1)
	ret0, _ := ret[0].(*auth.GetRefreshTokenResponseModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRefreshToken indicates an expected call of GetRefreshToken.
func (mr *MockStoragerMockRecorder) GetRefreshToken(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRefreshToken", reflect.TypeOf((*MockStorager)(nil).GetRefreshToken), arg0, arg1)
}

// GetUserByLogin mocks base method.
func (m *MockStorager) GetUserByLogin(arg0 context.Context, arg1 string) (*auth.GetUserByLoginResponseModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0, arg1)
	ret0, _ := ret[0].(*auth.GetUserByLoginResponseModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockStoragerMockRecorder) GetUserByLogin(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockStorager)(nil).GetUserByLogin), arg0, arg1)
}

// RemoveRefreshToken mocks base method.
func (m *MockStorager) RemoveRefreshToken(arg0 context.Context, arg1, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveRefreshToken", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveRefreshToken indicates an expected call of RemoveRefreshToken.
func (mr *MockStoragerMockRecorder) RemoveRefreshToken(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveRefreshToken", reflect.TypeOf((*MockStorager)(nil).RemoveRefreshToken), arg0, arg1, arg2)
}

// UpdateRefreshToken mocks base method.
func (m *MockStorager) UpdateRefreshToken(arg0 context.Context, arg1, arg2 string, arg3 time.Time) (*auth.UpdateRefreshTokenResponseModel, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRefreshToken", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(*auth.UpdateRefreshTokenResponseModel)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRefreshToken indicates an expected call of UpdateRefreshToken.
func (mr *MockStoragerMockRecorder) UpdateRefreshToken(arg0, arg1, arg2, arg3 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRefreshToken", reflect.TypeOf((*MockStorager)(nil).UpdateRefreshToken), arg0, arg1, arg2, arg3)
}
