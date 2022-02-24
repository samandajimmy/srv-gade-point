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

// NewReferralsHandler represent to register referrals endpoint
func NewReferralsHandler(echoGroup models.EchoGroup, us referrals.UseCase) {
	handler := &ReferralHandler{
		ReferralUseCase: us,
	}

	echoGroup.API.POST("/referral/core-trx", handler.coreTrx)
}

func (ref *ReferralHandler) coreTrx(echTx echo.Context) error {
	var referralCore models.CoreTrxPayload
	var responseData []models.CoreTrxResponse
	response = models.Response{}
	respErrors := &models.ResponseErrors{}
	err := echTx.Bind(&referralCore)

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	if err := echTx.Validate(referralCore); err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)

		return echTx.JSON(http.StatusBadRequest, response)
	}

	responseData, err = ref.ReferralUseCase.PostCoreTrx(echTx, referralCore)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)

		return echTx.JSON(getStatusCode(err), response)
	}

	response.SetResponse(responseData, respErrors)

	return echTx.JSON(getStatusCode(err), response)
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
