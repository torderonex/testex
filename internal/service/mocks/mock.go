// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	reflect "reflect"
	entities "testex/internal/entities"

	gomock "github.com/golang/mock/gomock"
)

// MockCommand is a mock of Command interface.
type MockCommand struct {
	ctrl     *gomock.Controller
	recorder *MockCommandMockRecorder
}

// MockCommandMockRecorder is the mock recorder for MockCommand.
type MockCommandMockRecorder struct {
	mock *MockCommand
}

// NewMockCommand creates a new mock instance.
func NewMockCommand(ctrl *gomock.Controller) *MockCommand {
	mock := &MockCommand{ctrl: ctrl}
	mock.recorder = &MockCommandMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCommand) EXPECT() *MockCommandMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCommand) Create(alias, script string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", alias, script)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCommandMockRecorder) Create(alias, script interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCommand)(nil).Create), alias, script)
}

// Execute mocks base method.
func (m *MockCommand) Execute(alias string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", alias)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockCommandMockRecorder) Execute(alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockCommand)(nil).Execute), alias)
}

// GetAll mocks base method.
func (m *MockCommand) GetAll() ([]entities.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].([]entities.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockCommandMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockCommand)(nil).GetAll))
}

// GetLogs mocks base method.
func (m *MockCommand) GetLogs(executedCommandId int) ([]entities.Log, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLogs", executedCommandId)
	ret0, _ := ret[0].([]entities.Log)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLogs indicates an expected call of GetLogs.
func (mr *MockCommandMockRecorder) GetLogs(executedCommandId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLogs", reflect.TypeOf((*MockCommand)(nil).GetLogs), executedCommandId)
}

// GetOne mocks base method.
func (m *MockCommand) GetOne(alias string) (entities.Command, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOne", alias)
	ret0, _ := ret[0].(entities.Command)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOne indicates an expected call of GetOne.
func (mr *MockCommandMockRecorder) GetOne(alias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOne", reflect.TypeOf((*MockCommand)(nil).GetOne), alias)
}

// StopCommand mocks base method.
func (m *MockCommand) StopCommand(id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StopCommand", id)
	ret0, _ := ret[0].(error)
	return ret0
}

// StopCommand indicates an expected call of StopCommand.
func (mr *MockCommandMockRecorder) StopCommand(id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StopCommand", reflect.TypeOf((*MockCommand)(nil).StopCommand), id)
}
