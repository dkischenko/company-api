// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock_company is a generated GoMock package.
package mock_company

import (
	reflect "reflect"

	models "githib.com/dkischenko/company-api/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockRepository) Create(company models.Company) (models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", company)
	ret0, _ := ret[0].(models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRepositoryMockRecorder) Create(company interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRepository)(nil).Create), company)
}

// CreateUser mocks base method.
func (m *MockRepository) CreateUser(user *models.User) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", user)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockRepositoryMockRecorder) CreateUser(user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockRepository)(nil).CreateUser), user)
}

// Delete mocks base method.
func (m *MockRepository) Delete(id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockRepositoryMockRecorder) Delete(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockRepository)(nil).Delete), id)
}

// FindOneUser mocks base method.
func (m *MockRepository) FindOneUser(name string) (models.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOneUser", name)
	ret0, _ := ret[0].(models.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOneUser indicates an expected call of FindOneUser.
func (mr *MockRepositoryMockRecorder) FindOneUser(name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOneUser", reflect.TypeOf((*MockRepository)(nil).FindOneUser), name)
}

// Get mocks base method.
func (m *MockRepository) Get(companyId uuid.UUID) (models.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", companyId)
	ret0, _ := ret[0].(models.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRepositoryMockRecorder) Get(companyId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRepository)(nil).Get), companyId)
}

// Update mocks base method.
func (m *MockRepository) Update(company *models.Company) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", company)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockRepositoryMockRecorder) Update(company interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockRepository)(nil).Update), company)
}