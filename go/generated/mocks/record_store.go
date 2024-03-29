// Code generated by MockGen. DO NOT EDIT.
// Source: storage/record_store.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	storage "github.com/openreserveio/dwn/go/storage"
)

// MockRecordStore is a mock of RecordStore interface.
type MockRecordStore struct {
	ctrl     *gomock.Controller
	recorder *MockRecordStoreMockRecorder
}

// MockRecordStoreMockRecorder is the mock recorder for MockRecordStore.
type MockRecordStoreMockRecorder struct {
	mock *MockRecordStore
}

// NewMockRecordStore creates a new mock instance.
func NewMockRecordStore(ctrl *gomock.Controller) *MockRecordStore {
	mock := &MockRecordStore{ctrl: ctrl}
	mock.recorder = &MockRecordStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRecordStore) EXPECT() *MockRecordStoreMockRecorder {
	return m.recorder
}

// AddMessageEntry mocks base method.
func (m *MockRecordStore) AddMessageEntry(entry *storage.MessageEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddMessageEntry", entry)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddMessageEntry indicates an expected call of AddMessageEntry.
func (mr *MockRecordStoreMockRecorder) AddMessageEntry(entry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddMessageEntry", reflect.TypeOf((*MockRecordStore)(nil).AddMessageEntry), entry)
}

// CreateRecord mocks base method.
func (m *MockRecordStore) CreateRecord(record *storage.Record, initialEntry *storage.MessageEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRecord", record, initialEntry)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateRecord indicates an expected call of CreateRecord.
func (mr *MockRecordStoreMockRecorder) CreateRecord(record, initialEntry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRecord", reflect.TypeOf((*MockRecordStore)(nil).CreateRecord), record, initialEntry)
}

// DeleteMessageEntry mocks base method.
func (m *MockRecordStore) DeleteMessageEntry(entry *storage.MessageEntry) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessageEntry", entry)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessageEntry indicates an expected call of DeleteMessageEntry.
func (mr *MockRecordStoreMockRecorder) DeleteMessageEntry(entry interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessageEntry", reflect.TypeOf((*MockRecordStore)(nil).DeleteMessageEntry), entry)
}

// DeleteMessageEntryByID mocks base method.
func (m *MockRecordStore) DeleteMessageEntryByID(messageEntryId string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteMessageEntryByID", messageEntryId)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteMessageEntryByID indicates an expected call of DeleteMessageEntryByID.
func (mr *MockRecordStoreMockRecorder) DeleteMessageEntryByID(messageEntryId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteMessageEntryByID", reflect.TypeOf((*MockRecordStore)(nil).DeleteMessageEntryByID), messageEntryId)
}

// GetMessageEntryByID mocks base method.
func (m *MockRecordStore) GetMessageEntryByID(messageEntryID string) *storage.MessageEntry {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMessageEntryByID", messageEntryID)
	ret0, _ := ret[0].(*storage.MessageEntry)
	return ret0
}

// GetMessageEntryByID indicates an expected call of GetMessageEntryByID.
func (mr *MockRecordStoreMockRecorder) GetMessageEntryByID(messageEntryID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMessageEntryByID", reflect.TypeOf((*MockRecordStore)(nil).GetMessageEntryByID), messageEntryID)
}

// GetRecord mocks base method.
func (m *MockRecordStore) GetRecord(recordId string) *storage.Record {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecord", recordId)
	ret0, _ := ret[0].(*storage.Record)
	return ret0
}

// GetRecord indicates an expected call of GetRecord.
func (mr *MockRecordStoreMockRecorder) GetRecord(recordId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecord", reflect.TypeOf((*MockRecordStore)(nil).GetRecord), recordId)
}

// GetRecordForCommit mocks base method.
func (m *MockRecordStore) GetRecordForCommit(recordId string) (*storage.Record, *storage.MessageEntry) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRecordForCommit", recordId)
	ret0, _ := ret[0].(*storage.Record)
	ret1, _ := ret[1].(*storage.MessageEntry)
	return ret0, ret1
}

// GetRecordForCommit indicates an expected call of GetRecordForCommit.
func (mr *MockRecordStoreMockRecorder) GetRecordForCommit(recordId interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRecordForCommit", reflect.TypeOf((*MockRecordStore)(nil).GetRecordForCommit), recordId)
}

// SaveRecord mocks base method.
func (m *MockRecordStore) SaveRecord(record *storage.Record) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveRecord", record)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveRecord indicates an expected call of SaveRecord.
func (mr *MockRecordStoreMockRecorder) SaveRecord(record interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveRecord", reflect.TypeOf((*MockRecordStore)(nil).SaveRecord), record)
}
