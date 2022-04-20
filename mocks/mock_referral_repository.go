// Code generated by MockGen. DO NOT EDIT.
// Source: gade/srv-gade-point/referrals (interfaces: RefRepository)

// Package mocks is a generated GoMock package.
package mocks

import (
	models "gade/srv-gade-point/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo"
)

// MockRefRepository is a mock of RefRepository interface.
type MockRefRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRefRepositoryMockRecorder
}

// MockRefRepositoryMockRecorder is the mock recorder for MockRefRepository.
type MockRefRepositoryMockRecorder struct {
	mock *MockRefRepository
}

// NewMockRefRepository creates a new mock instance.
func NewMockRefRepository(ctrl *gomock.Controller) *MockRefRepository {
	mock := &MockRefRepository{ctrl: ctrl}
	mock.recorder = &MockRefRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRefRepository) EXPECT() *MockRefRepositoryMockRecorder {
	return m.recorder
}

// RCreateReferral mocks base method.
func (m *MockRefRepository) RCreateReferral(arg0 echo.Context, arg1 models.ReferralCodes) (models.ReferralCodes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RCreateReferral", arg0, arg1)
	ret0, _ := ret[0].(models.ReferralCodes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RCreateReferral indicates an expected call of RCreateReferral.
func (mr *MockRefRepositoryMockRecorder) RCreateReferral(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RCreateReferral", reflect.TypeOf((*MockRefRepository)(nil).RCreateReferral), arg0, arg1)
}

// RFriendsReferral mocks base method.
func (m *MockRefRepository) RFriendsReferral(arg0 echo.Context, arg1 models.PayloadFriends) ([]models.Friends, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RFriendsReferral", arg0, arg1)
	ret0, _ := ret[0].([]models.Friends)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RFriendsReferral indicates an expected call of RFriendsReferral.
func (mr *MockRefRepositoryMockRecorder) RFriendsReferral(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RFriendsReferral", reflect.TypeOf((*MockRefRepository)(nil).RFriendsReferral), arg0, arg1)
}

// RGenerateCode mocks base method.
func (m *MockRefRepository) RGenerateCode(arg0 echo.Context, arg1 models.ReferralCodes, arg2 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RGenerateCode", arg0, arg1, arg2)
	ret0, _ := ret[0].(string)
	return ret0
}

// RGenerateCode indicates an expected call of RGenerateCode.
func (mr *MockRefRepositoryMockRecorder) RGenerateCode(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RGenerateCode", reflect.TypeOf((*MockRefRepository)(nil).RGenerateCode), arg0, arg1, arg2)
}

// RGetCampaignId mocks base method.
func (m *MockRefRepository) RGetCampaignId(arg0 echo.Context, arg1 string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RGetCampaignId", arg0, arg1)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RGetCampaignId indicates an expected call of RGetCampaignId.
func (mr *MockRefRepositoryMockRecorder) RGetCampaignId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RGetCampaignId", reflect.TypeOf((*MockRefRepository)(nil).RGetCampaignId), arg0, arg1)
}

// RGetHistoryIncentive mocks base method.
func (m *MockRefRepository) RGetHistoryIncentive(arg0 echo.Context, arg1 models.RequestHistoryIncentive) (models.ResponseHistoryIncentive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RGetHistoryIncentive", arg0, arg1)
	ret0, _ := ret[0].(models.ResponseHistoryIncentive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RGetHistoryIncentive indicates an expected call of RGetHistoryIncentive.
func (mr *MockRefRepositoryMockRecorder) RGetHistoryIncentive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RGetHistoryIncentive", reflect.TypeOf((*MockRefRepository)(nil).RGetHistoryIncentive), arg0, arg1)
}

// RGetReferralByCif mocks base method.
func (m *MockRefRepository) RGetReferralByCif(arg0 echo.Context, arg1 string) (models.ReferralCodes, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RGetReferralByCif", arg0, arg1)
	ret0, _ := ret[0].(models.ReferralCodes)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RGetReferralByCif indicates an expected call of RGetReferralByCif.
func (mr *MockRefRepositoryMockRecorder) RGetReferralByCif(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RGetReferralByCif", reflect.TypeOf((*MockRefRepository)(nil).RGetReferralByCif), arg0, arg1)
}

// RGetReferralCampaignMetadata mocks base method.
func (m *MockRefRepository) RGetReferralCampaignMetadata(arg0 echo.Context, arg1 models.PayloadValidator) (models.PrefixResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RGetReferralCampaignMetadata", arg0, arg1)
	ret0, _ := ret[0].(models.PrefixResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RGetReferralCampaignMetadata indicates an expected call of RGetReferralCampaignMetadata.
func (mr *MockRefRepositoryMockRecorder) RGetReferralCampaignMetadata(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RGetReferralCampaignMetadata", reflect.TypeOf((*MockRefRepository)(nil).RGetReferralCampaignMetadata), arg0, arg1)
}

// RGetReferralCodeByCampaignId mocks base method.
func (m *MockRefRepository) RGetReferralCodeByCampaignId(arg0 echo.Context, arg1 int64) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RGetReferralCodeByCampaignId", arg0, arg1)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RGetReferralCodeByCampaignId indicates an expected call of RGetReferralCodeByCampaignId.
func (mr *MockRefRepositoryMockRecorder) RGetReferralCodeByCampaignId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RGetReferralCodeByCampaignId", reflect.TypeOf((*MockRefRepository)(nil).RGetReferralCodeByCampaignId), arg0, arg1)
}

// RSumRefIncentive mocks base method.
func (m *MockRefRepository) RSumRefIncentive(arg0 echo.Context, arg1 string) (models.ObjIncentive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RSumRefIncentive", arg0, arg1)
	ret0, _ := ret[0].(models.ObjIncentive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RSumRefIncentive indicates an expected call of RSumRefIncentive.
func (mr *MockRefRepositoryMockRecorder) RSumRefIncentive(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RSumRefIncentive", reflect.TypeOf((*MockRefRepository)(nil).RSumRefIncentive), arg0, arg1)
}

// RTotalFriends mocks base method.
func (m *MockRefRepository) RTotalFriends(arg0 echo.Context, arg1 string) (models.RespTotalFriends, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RTotalFriends", arg0, arg1)
	ret0, _ := ret[0].(models.RespTotalFriends)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// RTotalFriends indicates an expected call of RTotalFriends.
func (mr *MockRefRepositoryMockRecorder) RTotalFriends(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RTotalFriends", reflect.TypeOf((*MockRefRepository)(nil).RTotalFriends), arg0, arg1)
}
