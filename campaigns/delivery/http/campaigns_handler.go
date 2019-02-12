package http

import (
	"context"
	"net/http"

	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"

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
	CUsecase campaigns.Usecase
}

// NewCampaignsHandler represent to register campaigns endpoint
func NewCampaignsHandler(e *echo.Echo, us campaigns.Usecase) {
	handler := &CampaignsHandler{
		CUsecase: us,
	}

	e.POST("/campaigns", handler.CreateCampaign)
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

	err = a.CUsecase.CreateCampaign(ctx, &campaign)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, campaign)
}

func isRequestValid(m *models.Campaign) (bool, error) {

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
