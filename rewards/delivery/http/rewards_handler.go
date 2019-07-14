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
	echoGroup.API.POST("/rewards/succeeded", handler.rewardSucceeded)
	echoGroup.API.POST("/rewards/rejected", handler.rewardRejected)
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
		response.Message = models.MessageNoRewards
	}

	requestLogger.Info("End of inquiry rewards.")

	return echTx.JSON(getStatusCode(err), response)
}

func (rwd *RewardHandler) rewardSucceeded(echTx echo.Context) error {
	var reward models.Reward
	response = models.Response{}
	err := echTx.Bind(&reward)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess

	return echTx.JSON(http.StatusOK, response)
}

func (rwd *RewardHandler) rewardRejected(echTx echo.Context) error {
	var reward models.Reward
	response = models.Response{}
	err := echTx.Bind(&reward)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess

	return echTx.JSON(http.StatusOK, response)
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
