package http

import (
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// CampaignsHandler represent the httphandler for campaigns
type CampaignsHandler struct {
	CampaignUseCase campaigns.UseCase
}

// NewCampaignsHandler represent to register campaigns endpoint
func NewCampaignsHandler(echoGroup models.EchoGroup, us campaigns.UseCase) {
	handler := &CampaignsHandler{
		CampaignUseCase: us,
	}

	//End Point For CMS
	echoGroup.Admin.POST("/campaigns", handler.CreateCampaign)
	echoGroup.Admin.PUT("/campaigns/status/:id", handler.UpdateStatusCampaign)
	echoGroup.Admin.GET("/campaigns", handler.GetCampaigns)
	echoGroup.Admin.GET("/campaigns/:id", handler.GetCampaignDetail)

	//End Point For External
	echoGroup.API.POST("/campaigns/value", handler.GetCampaignValue)
	echoGroup.API.GET("/campaigns/point", handler.GetUserPoint)
	echoGroup.API.GET("/campaigns/point/history", handler.GetUserPointHistory)
}

// CreateCampaign a handler to create a campaign
func (cmpgn *CampaignsHandler) CreateCampaign(c echo.Context) error {
	var campaign models.Campaign
	response = models.Response{}
	logger := models.RequestLogger{}
	err := c.Bind(&campaign)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	requestLogger := logger.GetRequestLogger(c, campaign)
	requestLogger.Info("Start to create a campaign.")
	err = cmpgn.CampaignUseCase.CreateCampaign(c, &campaign)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if (models.Campaign{}) != campaign {
		response.Data = campaign
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageSaveSuccess

	requestLogger.Info("End of create a campaign.")

	return c.JSON(http.StatusCreated, response)
}

// UpdateStatusCampaign a handler to update campaign status
func (cmpgn *CampaignsHandler) UpdateStatusCampaign(c echo.Context) error {
	updateCampaign := models.UpdateCampaign{}
	response = models.Response{}

	if err := c.Bind(&updateCampaign); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, updateCampaign)
	requestLogger.Info("Start to update a campaign.")
	err := cmpgn.CampaignUseCase.UpdateCampaign(c, c.Param("id"), &updateCampaign)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageUpdateSuccess
	requestLogger.Info("End of update a campaign.")

	return c.JSON(getStatusCode(err), response)
}

// GetCampaigns to get list of campaigns data
func (cmpgn *CampaignsHandler) GetCampaigns(c echo.Context) error {
	response = models.Response{}
	payload := map[string]interface{}{
		"name":      c.QueryParam("name"),
		"status":    c.QueryParam("status"),
		"startDate": c.QueryParam("startDate"),
		"endDate":   c.QueryParam("endDate"),
		"page":      c.QueryParam("page"),
		"limit":     c.QueryParam("limit"),
	}

	logger := models.RequestLogger{
		Payload: payload,
	}

	requestLogger := logger.GetRequestLogger(c, nil)
	requestLogger.Info("Start to get campaigns.")
	countCampaign, data, err := cmpgn.CampaignUseCase.GetCampaign(c, payload)

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
	response.TotalCount = countCampaign

	requestLogger.Info("End of get campaigns.")

	return c.JSON(http.StatusOK, response)
}

// GetCampaignDetail a handler  to provide and endpoint to get campaign detail
func (cmpgn *CampaignsHandler) GetCampaignDetail(c echo.Context) error {
	response = models.Response{}
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	requestLogger.Info("Start to get detail campaign.")

	responseData, err := cmpgn.CampaignUseCase.GetCampaignDetail(c, c.Param("id"))

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if (&models.Campaign{}) != responseData {
		response.Data = responseData
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	requestLogger.Info("End of get detail campaign.")

	return c.JSON(http.StatusOK, response)
}

// GetCampaignValue to validate point amount available and store the point trx
func (cmpgn *CampaignsHandler) GetCampaignValue(c echo.Context) error {
	campaignValue := models.GetCampaignValue{}
	response = models.Response{}

	if err := c.Bind(&campaignValue); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, campaignValue)
	requestLogger.Info("Start to get a campaign value.")
	userPoint, err := cmpgn.CampaignUseCase.GetCampaignValue(c, &campaignValue)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if (&models.UserPoint{}) != userPoint {
		response.Data = userPoint
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	requestLogger.Info("End of get a campaign value.")

	return c.JSON(http.StatusOK, response)
}

// GetUserPoint to get user current point
func (cmpgn *CampaignsHandler) GetUserPoint(c echo.Context) error {
	response = models.Response{}
	userID := c.QueryParam("userId")
	payload := map[string]interface{}{
		"userID": userID,
	}

	logger := models.RequestLogger{
		Payload: payload,
	}

	requestLogger := logger.GetRequestLogger(c, payload)
	requestLogger.Info("Start to get user point.")
	userPoint, _ := cmpgn.CampaignUseCase.GetUserPoint(c, userID)

	if (&models.UserPoint{}) != userPoint {
		response.Data = userPoint
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	requestLogger.Info("End of get user point.")

	return c.JSON(http.StatusOK, response)

}

// GetUserPointHistory is a handler to provide and endpoint to get user point history
func (cmpgn *CampaignsHandler) GetUserPointHistory(c echo.Context) error {
	response = models.Response{}
	userID := c.QueryParam("userId")
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
		"userID":      userID,
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
	data, counter, err := cmpgn.CampaignUseCase.GetUserPointHistory(c, payload)

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
