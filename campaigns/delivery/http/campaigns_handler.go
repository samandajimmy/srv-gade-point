package http

import (
	"context"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// Response represent the response
var response = models.Response{}

// CampaignsHandler represent the httphandler for campaigns
type CampaignsHandler struct {
	CampaignUseCase campaigns.UseCase
}

// NewCampaignsHandler represent to register campaigns endpoint
func NewCampaignsHandler(e *echo.Echo, us campaigns.UseCase) {
	handler := &CampaignsHandler{
		CampaignUseCase: us,
	}

	//End Point For CMS
	e.POST("/campaigns", handler.CreateCampaign)
	e.PUT("/campaigns/status/:id", handler.UpdateStatusCampaign)
	e.GET("/campaigns", handler.GetCampaigns)

	//End Point For External
	e.POST("/campaigns/value", handler.GetCampaignValue)
	e.GET("/campaigns/point", handler.GetUserPoint)
	e.GET("/campaigns/point/history", handler.GetUserPointHistory)
}

func (cmpgn *CampaignsHandler) CreateCampaign(c echo.Context) error {
	var campaign models.Campaign
	err := c.Bind(&campaign)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if ok, err := isRequestValid(&campaign); !ok {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusBadRequest, response)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = cmpgn.CampaignUseCase.CreateCampaign(ctx, &campaign)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MassageSaveSuccess
	response.Data = campaign
	return c.JSON(http.StatusCreated, response)
}

func (cmpgn *CampaignsHandler) UpdateStatusCampaign(c echo.Context) error {

	updateCampaign := new(models.UpdateCampaign)

	if err := c.Bind(updateCampaign); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if ok, err := isRequestValid(updateCampaign); !ok {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusBadRequest, response)
	}

	id, _ := strconv.Atoi(c.Param("id"))

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := cmpgn.CampaignUseCase.UpdateCampaign(ctx, int64(id), updateCampaign)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MassageUpdateSuccess
	response.Data = ""
	return c.JSON(http.StatusOK, response)
}

func (cmpgn *CampaignsHandler) GetCampaigns(c echo.Context) error {

	name := c.QueryParam("name")
	status := c.QueryParam("status")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := cmpgn.CampaignUseCase.GetCampaign(ctx, name, status, startDate, endDate)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.StatusSuccess
	response.Data = res
	return c.JSON(http.StatusOK, response)

}

func (cmpgn *CampaignsHandler) GetCampaignValue(c echo.Context) error {
	var campaignValue models.GetCampaignValue
	err := c.Bind(&campaignValue)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if ok, err := isRequestValid(&campaignValue); !ok {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusBadRequest, response)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	userPoint, err := cmpgn.CampaignUseCase.GetCampaignValue(ctx, &campaignValue)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MassagePointSuccess
	response.Data = userPoint
	return c.JSON(http.StatusOK, response)
}

func (cmpgn *CampaignsHandler) GetUserPoint(c echo.Context) error {

	userId := c.QueryParam("userId")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := cmpgn.CampaignUseCase.GetUserPoint(ctx, userId)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.StatusSuccess
	response.Data = res
	return c.JSON(http.StatusOK, response)

}

// GetUserPointHistory is a handler to provide and endpoint to get user point history
func (cmpgn *CampaignsHandler) GetUserPointHistory(c echo.Context) error {
	userID := c.QueryParam("userId")
	ctx := c.Request().Context()

	if ctx == nil {
		ctx = context.Background()
	}
	fmt.Println("cacing1")

	data, err := cmpgn.CampaignUseCase.GetUserPointHistory(ctx, userID)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.StatusSuccess
	response.Data = data

	// b, _ := json.Marshal(response)
	// fmt.Println(string(b))
	// err = json.UnmarshalJSON(b, &response)
	return c.JSON(http.StatusOK, response)
}

func isRequestValid(m interface{}) (bool, error) {

	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
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
