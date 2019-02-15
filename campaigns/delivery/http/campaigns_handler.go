package http

import (
	"context"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// CampaignsHandler represent the httphandler for campaigns
type CampaignsHandler struct {
	CampaignUseCase campaigns.UseCase
}

// NewCampaignsHandler represent to register campaigns endpoint
func NewCampaignsHandler(e *echo.Echo, us campaigns.UseCase) {
	handler := &CampaignsHandler{
		CampaignUseCase: us,
	}

	e.POST("/campaigns", handler.CreateCampaign)
	e.PUT("/statusCampaign/:id", handler.UpdateStatusCampaign)
	e.GET("/campaigns", handler.GetCampaigns)
}

func (a *CampaignsHandler) CreateCampaign(c echo.Context) error {
	var campaign models.Campaign
	err := c.Bind(&campaign)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&campaign); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.CampaignUseCase.CreateCampaign(ctx, &campaign)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, campaign)
}

func (a *CampaignsHandler) UpdateStatusCampaign(c echo.Context) error {

	updateCampaign := new(models.UpdateCampaign)

	if err := c.Bind(updateCampaign); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(updateCampaign); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	id, _ := strconv.Atoi(c.Param("id"))

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.CampaignUseCase.UpdateCampaign(ctx, int64(id), updateCampaign)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)
}

func (a *CampaignsHandler) GetCampaigns(c echo.Context) error {

	name := c.QueryParam("name")
	status := c.QueryParam("status")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.CampaignUseCase.GetCampaign(ctx, name, status, startDate, endDate)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusOK, res)

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
