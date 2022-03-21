// Code generated by MockGen. DO NOT EDIT.
// Source: gade/srv-gade-point/referrals (interfaces: RefUseCase)

// Package mocks is a generated GoMock package.
package mocks

import (
	models "gade/srv-gade-point/models"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	echo "github.com/labstack/echo"
)

// MockRefUseCase is a mock of RefUseCase interface.
type MockRefUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockRefUseCaseMockRecorder
}

// MockRefUseCaseMockRecorder is the mock recorder for MockRefUseCase.
type MockRefUseCaseMockRecorder struct {
	mock *MockRefUseCase
}

// NewMockRefUseCase creates a new mock instance.
func NewMockRefUseCase(ctrl *gomock.Controller) *MockRefUseCase {
	mock := &MockRefUseCase{ctrl: ctrl}
	mock.recorder = &MockRefUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRefUseCase) EXPECT() *MockRefUseCaseMockRecorder {
	return m.recorder
}

// UCreateReferralCodes mocks base method.
func (m *MockRefUseCase) UCreateReferralCodes(arg0 echo.Context, arg1 models.RequestCreateReferral) (models.RespReferral, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UCreateReferralCodes", arg0, arg1)
	ret0, _ := ret[0].(models.RespReferral)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UCreateReferralCodes indicates an expected call of UCreateReferralCodes.
func (mr *MockRefUseCaseMockRecorder) UCreateReferralCodes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UCreateReferralCodes", reflect.TypeOf((*MockRefUseCase)(nil).UCreateReferralCodes), arg0, arg1)
}

// UGetReferralCodes mocks base method.
func (m *MockRefUseCase) UGetReferralCodes(arg0 echo.Context, arg1 models.RequestReferralCodeUser) (models.ReferralCodeUser, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UGetReferralCodes", arg0, arg1)
	ret0, _ := ret[0].(models.ReferralCodeUser)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UGetReferralCodes indicates an expected call of UGetReferralCodes.
func (mr *MockRefUseCaseMockRecorder) UGetReferralCodes(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UGetReferralCodes", reflect.TypeOf((*MockRefUseCase)(nil).UGetReferralCodes), arg0, arg1)
}

// UGetPrefixActiveCampaignReferral mocks base method.
func (m *MockRefUseCase) UGetPrefixActiveCampaignReferral(arg0 echo.Context) (models.PrefixResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UGetPrefixActiveCampaignReferral", arg0)
	ret0, _ := ret[0].(models.PrefixResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UGetPrefixActiveCampaignReferral indicates an expected call of UGetPrefixActiveCampaignReferral.
func (mr *MockRefUseCaseMockRecorder) UGetPrefixActiveCampaignReferral(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UGetPrefixActiveCampaignReferral", reflect.TypeOf((*MockRefUseCase)(nil).UGetPrefixActiveCampaignReferral), arg0)
}

// UPostCoreTrx mocks base method.
func (m *MockRefUseCase) UPostCoreTrx(arg0 echo.Context, arg1 []models.CoreTrxPayload) ([]models.CoreTrxResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UPostCoreTrx", arg0, arg1)
	ret0, _ := ret[0].([]models.CoreTrxResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UPostCoreTrx indicates an expected call of UPostCoreTrx.
func (mr *MockRefUseCaseMockRecorder) UPostCoreTrx(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UPostCoreTrx", reflect.TypeOf((*MockRefUseCase)(nil).UPostCoreTrx), arg0, arg1)
}

// UValidateReferrer mocks base method.
func (m *MockRefUseCase) UValidateReferrer(arg0 echo.Context, arg1 models.PayloadValidator, arg2 *models.Campaign) (models.SumIncentive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UValidateReferrer", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.SumIncentive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UValidateReferrer indicates an expected call of UValidateReferrer.
func (mr *MockRefUseCaseMockRecorder) UValidateReferrer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UReferralCIFValidate", reflect.TypeOf((*MockRefUseCase)(nil).UReferralCIFValidate), arg0, arg1)
}

// UValidateReferrer mocks base method.
func (m *MockRefUseCase) UValidateReferrer(arg0 echo.Context, arg1 models.PayloadValidator, arg2 *models.Campaign) (models.SumIncentive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UValidateReferrer", arg0, arg1, arg2)
	ret0, _ := ret[0].(models.SumIncentive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UValidateReferrer indicates an expected call of UValidateReferrer.
func (mr *MockRefUseCaseMockRecorder) UValidateReferrer(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UValidateReferrer", reflect.TypeOf((*MockRefUseCase)(nil).UValidateReferrer), arg0, arg1, arg2)
}
