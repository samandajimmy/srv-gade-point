package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewardtrxs"
	"net/http"
	"regexp"
	"strconv"
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
func (rwd *RewardTrxHandler) getRewardTrxs(c echo.Context) error {
	response = models.Response{}
	idStr := c.QueryParam("id")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	if pageStr == "" {
		pageStr = "0"
	}

	if limitStr == "" {
		limitStr = "0"
	}

	payload := map[string]interface{}{
		"id":        idStr,
		"page":      pageStr,
		"limit":     limitStr,
		"startDate": startDate,
		"endDate":   endDate,
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

	dateFmtRgx := regexp.MustCompile(models.DateFormatRegex)

	if startDate != "" && !dateFmtRgx.MatchString(startDate) {
		requestLogger.Debug(models.ErrStartDateFormat)
		response.Status = models.StatusError
		response.Message = models.ErrStartDateFormat.Error()

		return c.JSON(http.StatusBadRequest, response)
	}

	if endDate != "" && !dateFmtRgx.MatchString(endDate) {
		requestLogger.Debug(models.ErrEndDateFormat)
		response.Status = models.StatusError
		response.Message = models.ErrEndDateFormat.Error()

		return c.JSON(http.StatusBadRequest, response)
	}

	payload["page"] = page
	payload["limit"] = limit
	requestLogger.Info("Start to get rewards transaction.")
	data, counter, err := rwd.RewardTrxUseCase.GetRewardTrxs(c, payload)

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

	requestLogger.Info("End of get reward transaction.")

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
