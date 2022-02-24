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
func NewReferralsHandler(echoGroup models.EchoGroup, rus referrals.UseCase) {
	handler := &ReferralHandler{
		ReferralUseCase: rus,
	}

	echoGroup.API.POST("/referral/generate", handler.GenerateReferralCodes)
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
