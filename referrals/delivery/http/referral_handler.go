package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// ReferralHandler represent the httphandler for referrals
type ReferralHandler struct {
	ReferralUseCase referrals.UseCase
}

// NewReferralsHandler represent to register referrals  endpoint
func NewReferralsHandler(echoGroup models.EchoGroup, us referrals.UseCase) {
	handler := &ReferralHandler{
		ReferralUseCase: us,
	}
	echoGroup.API.POST("/referral/generate", handler.GenerateReferralCodes)
	echoGroup.API.GET("/referral/get", handler.GetReferralCodes)
	echoGroup.API.POST("/referral/core-trx", handler.hcoreTrx)
}

func (ref *ReferralHandler) hcoreTrx(echTx echo.Context) error {
	var referralCore []models.CoreTrxPayload
	var responseData []models.CoreTrxResponse
	response = models.Response{}
	respErrors := &models.ResponseErrors{}
	err := echTx.Bind(&referralCore)

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	responseData, err = ref.ReferralUseCase.UPostCoreTrx(echTx, referralCore)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)

		return echTx.JSON(getStatusCode(err), response)
	}

	response.SetResponse(responseData, respErrors)

	return echTx.JSON(getStatusCode(err), response)
}

func (rc *ReferralHandler) GenerateReferralCodes(c echo.Context) error {
	var payload models.RequestCreateReferral

	respErrors := &models.ResponseErrors{}
	response = models.Response{}
	err := c.Bind(&payload)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if err = c.Validate(payload); err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)

		return c.JSON(http.StatusBadRequest, response)
	}

	data, err := rc.ReferralUseCase.CreateReferralCodes(c, payload)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = data
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageSaveSuccess
	response.Data = data

	return c.JSON(http.StatusCreated, response)
}

func (rc *ReferralHandler) GetReferralCodes(c echo.Context) error {
	var payload models.RequestReferralCodeUser

	respErrors := &models.ResponseErrors{}
	response = models.Response{}
	err := c.Bind(&payload)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if err = c.Validate(payload); err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)

		return c.JSON(http.StatusBadRequest, response)
	}

	data, err := rc.ReferralUseCase.GetReferralCodes(c, payload)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = data
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataFound
	response.Data = data

	return c.JSON(http.StatusFound, response)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if strings.Contains(err.Error(), "400") {
		return http.StatusBadRequest
	}

	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusOK
	}
}
