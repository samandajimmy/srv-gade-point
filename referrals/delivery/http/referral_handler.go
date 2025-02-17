package http

import (
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/middleware"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"

	"github.com/labstack/echo"
)

var (
	hCtrl = helper.NewHandler(&helper.Handler{})
)

// ReferralHandler represent the httphandler for referrals
type ReferralHandler struct {
	ReferralUseCase referrals.RefUseCase
}

// NewReferralsHandler represent to register referrals  endpoint
func NewReferralsHandler(echoGroup models.EchoGroup, us referrals.RefUseCase) {
	handler := &ReferralHandler{
		ReferralUseCase: us,
	}

	echoGroup.API.POST("/referral/generate", handler.HGenerateReferralCodes)
	middleware.ReferralAuth(handler.checkCif)
	echoGroup.API.POST("/referral/detail", handler.HGetReferralCodes)
	echoGroup.API.POST("/referral/prefix", handler.hGetPrefixCampaignReferral)
	echoGroup.API.POST("/referral/incentive", handler.HGetHistoriesIncentive)
	echoGroup.API.POST("/referral/total-friends", handler.HTotalFriends)
	echoGroup.API.POST("/referral/friends", handler.HFriendsReferral)
}

func (ref *ReferralHandler) HTotalFriends(c echo.Context) error {
	var pl models.RequestReferralCodeUser
	var errors models.ResponseErrors

	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := ref.ReferralUseCase.UTotalFriends(c, pl)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

func (ref *ReferralHandler) HGenerateReferralCodes(c echo.Context) error {
	var pl models.RequestCreateReferral
	var errors models.ResponseErrors
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := ref.ReferralUseCase.UCreateReferralCodes(c, pl)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

func (ref *ReferralHandler) HGetReferralCodes(c echo.Context) error {
	var pl models.RequestReferralCodeUser
	var errors models.ResponseErrors

	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := ref.ReferralUseCase.UGetReferralCodes(c, pl)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

func (ref *ReferralHandler) hGetPrefixCampaignReferral(c echo.Context) error {
	var errors models.ResponseErrors

	responseData, err := ref.ReferralUseCase.UGetPrefixActiveCampaignReferral(c)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	return hCtrl.ShowResponse(c, responseData, err, errors)
}
func (ref *ReferralHandler) HGetHistoriesIncentive(c echo.Context) error {
	var pl models.RequestHistoryIncentive
	var errors models.ResponseErrors

	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := ref.ReferralUseCase.UGetHistoryIncentive(c, pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

func (ref *ReferralHandler) checkCif(c echo.Context) error {
	var pl models.RequestReferralCodeUser
	var errors models.ResponseErrors
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := ref.ReferralUseCase.UReferralCIFValidate(c, pl.CIF)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

func (ref *ReferralHandler) HFriendsReferral(c echo.Context) error {
	var pl models.PayloadFriends
	var errors models.ResponseErrors
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := ref.ReferralUseCase.UFriendsReferral(c, pl)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}
