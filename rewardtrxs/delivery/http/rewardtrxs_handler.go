package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewardtrxs"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// RewardTrxHandler represent the httphandler for reward transaction
type RewardTrxHandler struct {
	RewardTrxUseCase rewardtrxs.UseCase
}

// NewRewardTrxHandler represent to register reward transaction endpoint
func NewRewardTrxHandler(echoGroup models.EchoGroup, us rewardtrxs.UseCase) {
	handler := &RewardTrxHandler{
		RewardTrxUseCase: us,
	}

	// End Point For CMS
	echoGroup.Admin.GET("/reward-transactions", handler.getRewardTrxs)

}

// getRewardTrxs a handler to get reward transaction
func (rwdTrx *RewardTrxHandler) getRewardTrxs(echTx echo.Context) error {
	var rwdPayload models.RewardsPayload
	response = models.Response{}
	respErrors := &models.ResponseErrors{}
	logger := models.RequestLogger{}
	err := echTx.Bind(&rwdPayload)
	logger.DataLog(echTx, rwdPayload).Info("Start to get rewards transaction.")

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)
		logger.DataLog(echTx, response).Info("End of get rewards transaction.")

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	if err := echTx.Validate(rwdPayload); err != nil {
		respErrors.SetTitle(err.Error())
		logger.DataLog(echTx, response).Info("End of get rewards transaction.")

		return echTx.JSON(http.StatusBadRequest, response)
	}

	dateFmtRgx := regexp.MustCompile(models.DateFormatRegex)

	if rwdPayload.StartTransactionDate != "" && !dateFmtRgx.MatchString(rwdPayload.StartTransactionDate) {
		respErrors.SetTitle(err.Error())
		logger.DataLog(echTx, response).Info("End of get rewards transaction.")

		return echTx.JSON(http.StatusBadRequest, response)
	}

	if rwdPayload.EndTransactionDate != "" && !dateFmtRgx.MatchString(rwdPayload.EndTransactionDate) {
		respErrors.SetTitle(err.Error())
		logger.DataLog(echTx, response).Info("End of get rewards transaction.")

		return echTx.JSON(http.StatusBadRequest, response)
	}

	if rwdPayload.StartSuccededDate != "" && !dateFmtRgx.MatchString(rwdPayload.StartSuccededDate) {
		respErrors.SetTitle(err.Error())
		logger.DataLog(echTx, response).Info("End of get rewards transaction.")

		return echTx.JSON(http.StatusBadRequest, response)
	}

	if rwdPayload.EndSuccededDate != "" && !dateFmtRgx.MatchString(rwdPayload.EndSuccededDate) {
		respErrors.SetTitle(err.Error())
		logger.DataLog(echTx, response).Info("End of get rewards transaction.")

		return echTx.JSON(http.StatusBadRequest, response)
	}

	responseData, counter, err := rwdTrx.RewardTrxUseCase.GetRewardTrxs(echTx, &rwdPayload)

	if err != nil {
		respErrors.SetTitle(err.Error())
		logger.DataLog(echTx, response).Info("End of get rewards transaction.")

		return echTx.JSON(http.StatusBadRequest, response)
	}

	if len(responseData) > 0 {
		response.Data = responseData
	}

	response.TotalCount = counter

	response.SetResponse(responseData, respErrors)
	logger.DataLog(echTx, response).Info("End of check rewards transaction.")

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
