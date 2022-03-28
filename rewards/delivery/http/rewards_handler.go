package http

import (
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewards"

	"github.com/labstack/echo"
)

var (
	hCtrl = helper.NewHandler(&helper.Handler{})
)

// RewardHandler represent the httphandler for reward
type RewardHandler struct {
	RewardUseCase rewards.UseCase
}

// NewRewardHandler represent to register reward endpoint
func NewRewardHandler(echoGroup models.EchoGroup, us rewards.UseCase) {
	handler := &RewardHandler{
		RewardUseCase: us,
	}

	// End Point For CMS
	echoGroup.Admin.GET("/rewards", handler.getRewards)

	// End Point For External
	echoGroup.API.POST("/rewards/inquiry", handler.rewardInquiry)
	echoGroup.API.POST("/rewards/inquiry/multi", handler.MultiRewardInquiry)
	echoGroup.API.POST("/rewards/succeeded", handler.rewardPayment)
	echoGroup.API.POST("/rewards/rejected", handler.rewardPayment)
	echoGroup.API.POST("/rewards/check-transaction", handler.checkTransaction)
	echoGroup.API.GET("/reward/promotions", handler.getRewardPromotions)

	// Endpoint for
	echoGroup.API.POST("/rewards/core-trx", handler.hCoreTrx)
}

func (rwd *RewardHandler) hCoreTrx(c echo.Context) error {
	var pl []models.CoreTrxPayload
	var errors models.ResponseErrors

	if err := c.Bind(&pl); err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := rwd.RewardUseCase.UPostCoreTrx(c, pl)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

func (rwd *RewardHandler) MultiRewardInquiry(c echo.Context) error {
	var pl models.PayloadValidator
	var errors = &models.ResponseErrors{}
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, *errors)
	}

	pl.IsMulti = true
	responseData, errors := rwd.RewardUseCase.Inquiry(c, &pl)

	return hCtrl.ShowResponse(c, responseData, nil, *errors)
}

func (rwd *RewardHandler) rewardInquiry(c echo.Context) error {
	var pl models.PayloadValidator
	var errors = &models.ResponseErrors{}
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, *errors)
	}

	responseData, errors := rwd.RewardUseCase.Inquiry(c, &pl)

	return hCtrl.ShowResponse(c, responseData, nil, *errors)
}

func (rwd *RewardHandler) rewardPayment(c echo.Context) error {
	var pl models.RewardPayment
	var errors models.ResponseErrors
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := rwd.RewardUseCase.Payment(c, &pl)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

func (rwd *RewardHandler) checkTransaction(c echo.Context) error {
	var pl models.RewardPayment
	var errors models.ResponseErrors
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, err := rwd.RewardUseCase.Payment(c, &pl)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

// getRewards a handler to get rewards
func (rwd *RewardHandler) getRewards(c echo.Context) error {
	var pl models.RewardsPayload
	var errors models.ResponseErrors
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, errors)
	}

	responseData, counter, err := rwd.RewardUseCase.GetRewards(c, &pl)
	hCtrl.SetTotalCount(counter)

	return hCtrl.ShowResponse(c, responseData, err, errors)
}

// getRewardPromotions Get reward promotions by param promoCode, transactionDate and transactionAmount
func (rwd *RewardHandler) getRewardPromotions(c echo.Context) error {
	var pl models.RewardPromotionLists
	var errors = &models.ResponseErrors{}
	err := hCtrl.Validate(c, &pl)

	if err != nil {
		return hCtrl.ShowResponse(c, nil, err, *errors)
	}

	responseData, errors, err := rwd.RewardUseCase.GetRewardPromotions(c, pl)

	return hCtrl.ShowResponse(c, responseData, err, *errors)
}
