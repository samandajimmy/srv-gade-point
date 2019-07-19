package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewards"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// RewardHandler represent the httphandler for reward
type RewardHandler struct {
	RewardUseCase rewards.UseCase
}

// NewRewardHandler represent to register reward endpoint
func NewRewardHandler(echoGroup models.EchoGroup, us rewards.UseCase) {
	handler := &RewardHandler{
		RewardUseCase: us,
	}

	// Dummy End Point
	echoGroup.API.POST("/rewards/inquiry", handler.rewardInquiry)
	echoGroup.API.POST("/rewards/succeeded", handler.rewardPayment)
	echoGroup.API.POST("/rewards/rejected", handler.rewardPayment)
}

func (rwd *RewardHandler) rewardInquiry(echTx echo.Context) error {
	var plValidator models.PayloadValidator
	respErrors := &models.ResponseErrors{}
	response = models.Response{}
	logger := models.RequestLogger{}
	err := echTx.Bind(&plValidator)
	logger.DataLog(echTx, plValidator).Info("Start to inquiry rewards.")

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)
		logger.DataLog(echTx, response).Info("End of inquiry rewards.")

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	if err = echTx.Validate(plValidator); err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(echTx, response).Info("End of inquiry rewards.")

		return echTx.JSON(http.StatusBadRequest, response)
	}
	responseData, respErrors := rwd.RewardUseCase.Inquiry(echTx, &plValidator)
	response.SetResponse(responseData, respErrors)
	logger.DataLog(echTx, response).Info("End of inquiry rewards.")

	return echTx.JSON(getStatusCode(err), response)
}

func (rwd *RewardHandler) rewardPayment(echTx echo.Context) error {
	var rwdPayment models.RewardPayment
	response = models.Response{}
	respErrors := &models.ResponseErrors{}
	logger := models.RequestLogger{}
	err := echTx.Bind(&rwdPayment)
	logger.DataLog(echTx, rwdPayment).Info("Start to payment rewards.")

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)
		logger.DataLog(echTx, response).Info("End of payment rewards.")

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	if err := echTx.Validate(rwdPayment); err != nil {
		respErrors.SetTitle(err.Error())
		logger.DataLog(echTx, response).Info("End of payment rewards.")

		return echTx.JSON(http.StatusBadRequest, response)
	}

	err = rwd.RewardUseCase.Payment(echTx, &rwdPayment)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(echTx, response).Info("End of payment rewards.")

		return echTx.JSON(getStatusCode(err), response)
	}

	response.SetResponse(nil, respErrors)
	logger.DataLog(echTx, response).Info("End of payment rewards.")

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
