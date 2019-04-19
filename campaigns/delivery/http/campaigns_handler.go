package http

import (
	"context"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"net/http"
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
	name := c.QueryParam("name")
	status := c.QueryParam("status")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	ctx := c.Request().Context()
	response.Data = ""
	response.TotalCount = ""

	if ctx == nil {
		ctx = context.Background()
	}

	countCampaign, res, err := cmpgn.CampaignUseCase.GetCampaign(ctx, name, status, startDate, endDate, int(page), int(limit))

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if len(res) > 0 {
		response.Data = res
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.TotalCount = countCampaign

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
	var campaignValue models.GetCampaignValue
	response = models.Response{}
	ctx := c.Request().Context()
	err := c.Bind(&campaignValue)

	if ctx == nil {
		ctx = context.Background()
	}

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	userPoint, err := cmpgn.CampaignUseCase.GetCampaignValue(ctx, &campaignValue)

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
	return c.JSON(http.StatusOK, response)
}

// GetUserPoint to get user current point
func (cmpgn *CampaignsHandler) GetUserPoint(c echo.Context) error {
	response = models.Response{}
	userID := c.QueryParam("userId")
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	userPoint, err := cmpgn.CampaignUseCase.GetUserPoint(ctx, userID)

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
	return c.JSON(http.StatusOK, response)

}

// GetUserPointHistory is a handler to provide and endpoint to get user point history
func (cmpgn *CampaignsHandler) GetUserPointHistory(c echo.Context) error {
	response = models.Response{}
	userID := c.QueryParam("userId")
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}

	data, err := cmpgn.CampaignUseCase.GetUserPointHistory(ctx, userID)

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
