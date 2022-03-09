// Code generated by MockGen. DO NOT EDIT.
// Source: gade/srv-gade-point/rewardtrxs (interfaces: RtRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	models "gade/srv-gade-point/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo"
)

// MockRtRepository is a mock of RtRepository interface.
type MockRtRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRtRepositoryMockRecorder
}

// MockRtRepositoryMockRecorder is the mock recorder for MockRtRepository.
type MockRtRepositoryMockRecorder struct {
	mock *MockRtRepository
}

// NewMockRtRepository creates a new mock instance.
func NewMockRtRepository(ctrl *gomock.Controller) *MockRtRepository {
	mock := &MockRtRepository{ctrl: ctrl}
	mock.recorder = &MockRtRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRtRepository) EXPECT() *MockRtRepositoryMockRecorder {
	return m.recorder
}

// CheckRefID mocks base method.
func (m *MockRtRepository) CheckRefID(arg0 echo.Context, arg1 string) (*models.RewardTrx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRefID", arg0, arg1)
	ret0, _ := ret[0].(*models.RewardTrx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRefID indicates an expected call of CheckRefID.
func (mr *MockRtRepositoryMockRecorder) CheckRefID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRefID", reflect.TypeOf((*MockRtRepository)(nil).CheckRefID), arg0, arg1)
}

// CheckRootRefId mocks base method.
func (m *MockRtRepository) CheckRootRefId(arg0 echo.Context, arg1 string) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckRootRefId", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckRootRefId indicates an expected call of CheckRootRefId.
func (mr *MockRtRepositoryMockRecorder) CheckRootRefId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckRootRefId", reflect.TypeOf((*MockRtRepository)(nil).CheckRootRefId), arg0, arg1)
}

// CheckTrx mocks base method.
func (m *MockRtRepository) CheckTrx(arg0 echo.Context, arg1 string) (*models.RewardTrx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckTrx", arg0, arg1)
	ret0, _ := ret[0].(*models.RewardTrx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckTrx indicates an expected call of CheckTrx.
func (mr *MockRtRepositoryMockRecorder) CheckTrx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckTrx", reflect.TypeOf((*MockRtRepository)(nil).CheckTrx), arg0, arg1)
}

// CountByCIF mocks base method.
func (m *MockRtRepository) CountByCIF(arg0 echo.Context, arg1 models.Quota, arg2 models.Reward, arg3 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountByCIF", arg0, arg1, arg2, arg3)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountByCIF indicates an expected call of CountByCIF.
func (mr *MockRtRepositoryMockRecorder) CountByCIF(arg0, arg1, arg2, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountByCIF", reflect.TypeOf((*MockRtRepository)(nil).CountByCIF), arg0, arg1, arg2, arg3)
}

// CountRewardTrxs mocks base method.
func (m *MockRtRepository) CountRewardTrxs(arg0 echo.Context, arg1 *models.RewardsPayload) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CountRewardTrxs", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CountRewardTrxs indicates an expected call of CountRewardTrxs.
func (mr *MockRtRepositoryMockRecorder) CountRewardTrxs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CountRewardTrxs", reflect.TypeOf((*MockRtRepository)(nil).CountRewardTrxs), arg0, arg1)
}

// Create mocks base method.
func (m *MockRtRepository) Create(arg0 echo.Context, arg1 models.PayloadValidator, arg2 models.RewardsInquiry) ([]*models.RewardTrx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*models.RewardTrx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockRtRepositoryMockRecorder) Create(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockRtRepository)(nil).Create), arg0, arg1, arg2)
}

// GetByRefID mocks base method.
func (m *MockRtRepository) GetByRefID(arg0 echo.Context, arg1 string) (models.RewardsInquiry, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRefID", arg0, arg1)
	ret0, _ := ret[0].(models.RewardsInquiry)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByRefID indicates an expected call of GetByRefID.
func (mr *MockRtRepositoryMockRecorder) GetByRefID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRefID", reflect.TypeOf((*MockRtRepository)(nil).GetByRefID), arg0, arg1)
}

// GetInquiredTrx mocks base method.
func (m *MockRtRepository) GetInquiredTrx() ([]models.RewardTrx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetInquiredTrx")
	ret0, _ := ret[0].([]models.RewardTrx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetInquiredTrx indicates an expected call of GetInquiredTrx.
func (mr *MockRtRepositoryMockRecorder) GetInquiredTrx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetInquiredTrx", reflect.TypeOf((*MockRtRepository)(nil).GetInquiredTrx))
}

// GetRewardByPayload mocks base method.
func (m *MockRtRepository) GetRewardByPayload(arg0 echo.Context, arg1 models.PayloadValidator, arg2 *models.VoucherCode) ([]*models.Reward, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRewardByPayload", arg0, arg1, arg2)
	ret0, _ := ret[0].([]*models.Reward)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRewardByPayload indicates an expected call of GetRewardByPayload.
func (mr *MockRtRepositoryMockRecorder) GetRewardByPayload(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRewardByPayload", reflect.TypeOf((*MockRtRepository)(nil).GetRewardByPayload), arg0, arg1, arg2)
}

// GetRewardTrxs mocks base method.
func (m *MockRtRepository) GetRewardTrxs(arg0 echo.Context, arg1 *models.RewardsPayload) ([]models.RewardTrx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRewardTrxs", arg0, arg1)
	ret0, _ := ret[0].([]models.RewardTrx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRewardTrxs indicates an expected call of GetRewardTrxs.
func (mr *MockRtRepositoryMockRecorder) GetRewardTrxs(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRewardTrxs", reflect.TypeOf((*MockRtRepository)(nil).GetRewardTrxs), arg0, arg1)
}

// RewardTrxTimeout mocks base method.
func (m *MockRtRepository) RewardTrxTimeout(arg0 models.RewardTrx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RewardTrxTimeout", arg0)
}

// RewardTrxTimeout indicates an expected call of RewardTrxTimeout.
func (mr *MockRtRepositoryMockRecorder) RewardTrxTimeout(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RewardTrxTimeout", reflect.TypeOf((*MockRtRepository)(nil).RewardTrxTimeout), arg0)
}

// UpdateRewardTrx mocks base method.
func (m *MockRtRepository) UpdateRewardTrx(arg0 echo.Context, arg1 *models.RewardPayment, arg2 int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRewardTrx", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRewardTrx indicates an expected call of UpdateRewardTrx.
func (mr *MockRtRepositoryMockRecorder) UpdateRewardTrx(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRewardTrx", reflect.TypeOf((*MockRtRepository)(nil).UpdateRewardTrx), arg0, arg1, arg2)
}

// UpdateTimeoutTrx mocks base method.
func (m *MockRtRepository) UpdateTimeoutTrx() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTimeoutTrx")
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTimeoutTrx indicates an expected call of UpdateTimeoutTrx.
func (mr *MockRtRepositoryMockRecorder) UpdateTimeoutTrx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTimeoutTrx", reflect.TypeOf((*MockRtRepository)(nil).UpdateTimeoutTrx))
}
