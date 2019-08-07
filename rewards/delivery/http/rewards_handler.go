package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewards"
	"net/http"
	"strconv"
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

	// End Point For CMS
	echoGroup.Admin.GET("/rewards", handler.getRewards)

	// End Point For External
	echoGroup.API.POST("/rewards/inquiry", handler.rewardInquiry)
	echoGroup.API.POST("/rewards/succeeded", handler.rewardPayment)
	echoGroup.API.POST("/rewards/rejected", handler.rewardPayment)
	echoGroup.API.POST("/rewards/check-transaction", handler.checkTransaction)
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
	var responseData models.RewardTrxResponse
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

	responseData, err = rwd.RewardUseCase.Payment(echTx, &rwdPayment)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(echTx, response).Info("End of payment rewards.")

		return echTx.JSON(getStatusCode(err), response)
	}

	if responseData.Status != "" {
		respErrors.SetTitle(models.ErrRefIDStatus.Error() + responseData.Status)
	}

	response.SetResponse(nil, respErrors)

	logger.DataLog(echTx, response).Info("End of payment rewards.")

	return echTx.JSON(getStatusCode(err), response)
}

func (rwd *RewardHandler) checkTransaction(echTx echo.Context) error {
	var rwdPayment models.RewardPayment
	var responseData models.RewardTrxResponse
	response = models.Response{}
	respErrors := &models.ResponseErrors{}
	logger := models.RequestLogger{}
	err := echTx.Bind(&rwdPayment)
	logger.DataLog(echTx, rwdPayment).Info("Start to check rewards transaction.")

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)
		logger.DataLog(echTx, response).Info("End of check rewards transaction.")

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	if err := echTx.Validate(rwdPayment); err != nil {
		respErrors.SetTitle(err.Error())
		logger.DataLog(echTx, response).Info("End of check rewards transaction.")

		return echTx.JSON(http.StatusBadRequest, response)
	}

	responseData, err = rwd.RewardUseCase.CheckTransaction(echTx, &rwdPayment)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(echTx, response).Info("End of check rewards transaction.")

		return echTx.JSON(getStatusCode(err), response)
	}

	if models.RewardTrxInquired != *responseData.StatusCode {
		respErrors.SetTitle(models.ErrRefIDStatus.Error() + responseData.Status)
	}

	response.SetResponse(nil, respErrors)
	logger.DataLog(echTx, response).Info("End of check rewards transaction.")

	if models.RewardTrxInquired == *responseData.StatusCode {
		response.Message = models.ErrRefIDStatus.Error() + responseData.Status
	}

	return echTx.JSON(getStatusCode(err), response)
}

// getRewards a handler to get rewards
func (rwd *RewardHandler) getRewards(c echo.Context) error {
	response = models.Response{}
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	if pageStr == "" {
		pageStr = "0"
	}

	if limitStr == "" {
		limitStr = "0"
	}

	payload := map[string]interface{}{
		"page":  pageStr,
		"limit": limitStr,
	}

	logger := models.RequestLogger{
		Payload: payload,
	}

	requestLogger := logger.GetRequestLogger(c, nil)

	page, err := strconv.Atoi(payload["page"].(string))

	if err != nil {
		requestLogger.Debug(err)
		response.Status = models.StatusError
		response.Message = http.StatusText(http.StatusBadRequest)

		return c.JSON(http.StatusBadRequest, response)
	}

	limit, err := strconv.Atoi(payload["limit"].(string))

	if err != nil {
		requestLogger.Debug(err)
		response.Status = models.StatusError
		response.Message = http.StatusText(http.StatusBadRequest)

		return c.JSON(http.StatusBadRequest, response)
	}

	payload["page"] = page
	payload["limit"] = limit
	requestLogger.Info("Start to get rewards.")
	data, counter, err := rwd.RewardUseCase.GetRewards(c, payload)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if len(data) > 0 {
		response.Data = data
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.TotalCount = counter

	requestLogger.Info("End of get rewards.")

	return c.JSON(getStatusCode(err), response)
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
