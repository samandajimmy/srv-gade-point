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
	echoGroup.API.POST("/referral/core-trx", handler.hCoreTrx)
	middleware.ReferralAuth(handler.checkCif)
	echoGroup.API.GET("/referral/detail", handler.HGetReferralCodes)
	echoGroup.API.GET("/referral/prefix", handler.hGetPrefixCampaignReferral)
}

func (ref *ReferralHandler) hCoreTrx(c echo.Context) error {
	var pl []models.CoreTrxPayload
	var errors models.ResponseErrors

	if err := c.Bind(&pl); err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := ref.ReferralUseCase.UPostCoreTrx(c, pl)

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
