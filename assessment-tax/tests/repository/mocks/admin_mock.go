// Code generated by MockGen. DO NOT EDIT.
// Source: ../../repository/admin.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	model "github.com/LGROW101/assessment-tax/model"
	gomock "github.com/golang/mock/gomock"
)

// MockAdminRepository is a mock of AdminRepository interface.
type MockAdminRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAdminRepositoryMockRecorder
}

// MockAdminRepositoryMockRecorder is the mock recorder for MockAdminRepository.
type MockAdminRepositoryMockRecorder struct {
	mock *MockAdminRepository
}

// NewMockAdminRepository creates a new mock instance.
func NewMockAdminRepository(ctrl *gomock.Controller) *MockAdminRepository {
	mock := &MockAdminRepository{ctrl: ctrl}
	mock.recorder = &MockAdminRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAdminRepository) EXPECT() *MockAdminRepositoryMockRecorder {
	return m.recorder
}

// GetConfig mocks base method.
func (m *MockAdminRepository) GetConfig() (*model.AdminConfig, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConfig")
	ret0, _ := ret[0].(*model.AdminConfig)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetConfig indicates an expected call of GetConfig.
func (mr *MockAdminRepositoryMockRecorder) GetConfig() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConfig", reflect.TypeOf((*MockAdminRepository)(nil).GetConfig))
}

// InsertConfig mocks base method.
func (m *MockAdminRepository) InsertConfig(config *model.AdminConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertConfig", config)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertConfig indicates an expected call of InsertConfig.
func (mr *MockAdminRepositoryMockRecorder) InsertConfig(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertConfig", reflect.TypeOf((*MockAdminRepository)(nil).InsertConfig), config)
}

// UpdateConfig mocks base method.
func (m *MockAdminRepository) UpdateConfig(config *model.AdminConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateConfig", config)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateConfig indicates an expected call of UpdateConfig.
func (mr *MockAdminRepositoryMockRecorder) UpdateConfig(config interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateConfig", reflect.TypeOf((*MockAdminRepository)(nil).UpdateConfig), config)
}
