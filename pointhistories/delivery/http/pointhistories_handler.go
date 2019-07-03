package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/pointhistories"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// PointHistoriesHandler represent the httphandler for pointhistories
type PointHistoriesHandler struct {
	PointHistoryUseCase pointhistories.UseCase
}

// NewPointHistoriesHandler represent to register pointhistories endpoint
func NewPointHistoriesHandler(echoGroup models.EchoGroup, pointHistUs pointhistories.UseCase) {
	handler := &PointHistoriesHandler{PointHistoryUseCase: pointHistUs}

	// End Point For CMS
	echoGroup.Admin.POST("/point", handler.GetUserPoint)
	echoGroup.Admin.POST("/point/history", handler.GetUserPointHistory)
	echoGroup.Admin.POST("/point/histories", handler.GetUsers)

	// End Point For External
	echoGroup.API.POST("/point", handler.GetUserPoint)
	echoGroup.API.POST("/point/history", handler.GetUserPointHistory)
}

// GetUsers a handler to create a pointHistory
func (pointHist *PointHistoriesHandler) GetUsers(c echo.Context) error {
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
	requestLogger.Info("Start to get pointHistory users.")
	data, counter, err := pointHist.PointHistoryUseCase.GetUsers(c, payload)

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

	requestLogger.Info("End of get pointHistory users.")

	return c.JSON(getStatusCode(err), response)
}

// GetUserPointHistory is a handler to provide and endpoint to get user point history
func (pointHist *PointHistoriesHandler) GetUserPointHistory(c echo.Context) error {
	response = models.Response{}
	CIF := c.QueryParam("CIF")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")
	startDateRg := c.QueryParam("startDateRg")
	endDateRg := c.QueryParam("endDateRg")

	// validate page and limit string input
	if pageStr == "" {
		pageStr = "0"
	}

	if limitStr == "" {
		limitStr = "0"
	}

	// prepare payload for logger
	payload := map[string]interface{}{
		"CIF":         CIF,
		"page":        pageStr,
		"limit":       limitStr,
		"startDateRg": startDateRg,
		"endDateRg":   endDateRg,
	}

	logger := models.RequestLogger{
		Payload: payload,
	}

	requestLogger := logger.GetRequestLogger(c, payload)
	requestLogger.Info("Start to get user point history")

	// validate payload values
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

	if startDateRg != "" && !dateFmtRgx.MatchString(startDateRg) {
		requestLogger.Debug(models.ErrStartDateFormat)
		response.Status = models.StatusError
		response.Message = models.ErrStartDateFormat.Error()

		return c.JSON(http.StatusBadRequest, response)
	}

	if endDateRg != "" && !dateFmtRgx.MatchString(endDateRg) {
		requestLogger.Debug(models.ErrEndDateFormat)
		response.Status = models.StatusError
		response.Message = models.ErrEndDateFormat.Error()

		return c.JSON(http.StatusBadRequest, response)
	}

	payload["page"] = page
	payload["limit"] = limit
	data, counter, err := pointHist.PointHistoryUseCase.GetUserPointHistory(c, payload)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if len(data) > 0 {
		response.Data = data
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.TotalCount = counter
	requestLogger.Info("End of get user point history")

	return c.JSON(http.StatusOK, response)
}

// GetUserPoint to get user current point
func (pointHist *PointHistoriesHandler) GetUserPoint(c echo.Context) error {
	response = models.Response{}
	CIF := c.QueryParam("CIF")
	payload := map[string]interface{}{
		"CIF": CIF,
	}

	logger := models.RequestLogger{
		Payload: payload,
	}

	requestLogger := logger.GetRequestLogger(c, payload)
	requestLogger.Info("Start to get user point.")
	userPoint, _ := pointHist.PointHistoryUseCase.GetUserPoint(c, CIF)

	if (&models.UserPoint{}) != userPoint {
		response.Data = userPoint
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	requestLogger.Info("End of get user point.")

	return c.JSON(http.StatusOK, response)

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
		return http.StatusInternalServerError
	}
}
