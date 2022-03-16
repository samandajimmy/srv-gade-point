package http

import (
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/middleware"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"

	"github.com/labstack/echo"
)

var (
	hCtrl helper.IHandler
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

	hCtrl = helper.NewHandler(&helper.Handler{})
	echoGroup.API.POST("/referral/generate", handler.hGenerateReferralCodes)
	echoGroup.API.POST("/referral/core-trx", handler.hCoreTrx)
	middleware.ReferralAuth(handler.checkCif)
	echoGroup.API.GET("/referral/get", handler.hGetReferralCodes)
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

func (ref *ReferralHandler) hGenerateReferralCodes(c echo.Context) error {
	var pl models.RequestCreateReferral
	var errors models.ResponseErrors
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := ref.ReferralUseCase.CreateReferralCodes(c, pl)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

func (ref *ReferralHandler) hGetReferralCodes(c echo.Context) error {
	var pl models.RequestReferralCodeUser
	var errors models.ResponseErrors
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := ref.ReferralUseCase.GetReferralCodes(c, pl)

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
