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

	echoGroup.API.POST("/referral/generate", handler.CreateReferralCodes)
}

func (rc *ReferralHandler) CreateReferralCodes(c echo.Context) error {
	var referralcodes models.ReferralCodes

	respErrors := &models.ResponseErrors{}
	response = models.Response{}
	logger := models.RequestLogger{}
	err := c.Bind(&referralcodes)

	requestLogger := logger.GetRequestLogger(c, referralcodes)
	requestLogger.Info("Start to create a referral code.")

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if err = c.Validate(referralcodes); err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End to create referral code for client")

		return c.JSON(http.StatusBadRequest, response)
	}

	err = rc.ReferralUseCase.CreateReferralCodes(c, &referralcodes)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageSaveSuccess
	response.Data = referralcodes

	requestLogger.Info("End of create a referral code.")

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
