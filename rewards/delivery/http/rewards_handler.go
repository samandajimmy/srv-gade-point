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
	response = models.Response{}
	err := echTx.Bind(&plValidator)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	if err = echTx.Validate(plValidator); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(http.StatusBadRequest, response)
	}

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(echTx, plValidator)
	requestLogger.Info("Start to inquiry rewards.")
	responseData, err := rwd.RewardUseCase.Inquiry(echTx, &plValidator)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess

	if (models.RewardsInquiry{}) != responseData {
		response.Data = responseData
	} else {
		response.Status = models.StatusError
		response.Message = models.MessageNoRewards
	}

	requestLogger.Info("End of inquiry rewards.")

	return echTx.JSON(getStatusCode(err), response)
}

func (rwd *RewardHandler) rewardPayment(echTx echo.Context) error {
	var rwdPayment models.RewardPayment
	response = models.Response{}

	if err := echTx.Bind(&rwdPayment); err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	if err := echTx.Validate(rwdPayment); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(http.StatusBadRequest, response)
	}

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(echTx, rwdPayment)
	requestLogger.Info("Start to payment rewards.")
	err := rwd.RewardUseCase.Payment(echTx, &rwdPayment)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess

	requestLogger.Info("End of payment rewards.")

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
