package http

import (
	"gade/srv-gade-point/campaigns"
	_campaignHttpDelivery "gade/srv-gade-point/campaigns/delivery/http"
	"gade/srv-gade-point/campaigntrxs"
	"gade/srv-gade-point/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// CampaigntrxsHandler represent the httphandler for campaigntrxs
type CampaigntrxsHandler struct {
	CampaigntrxUseCase campaigntrxs.UseCase
}

// NewCampaignTrxsHandler represent to register campaigntrxs endpoint
func NewCampaignTrxsHandler(echoGroup models.EchoGroup, cmpTrxUs campaigntrxs.UseCase, cmpUs campaigns.UseCase) {
	handler := &CampaigntrxsHandler{CampaigntrxUseCase: cmpTrxUs}
	cmpHandler := &_campaignHttpDelivery.CampaignsHandler{CampaignUseCase: cmpUs}

	// End Point For CMS
	echoGroup.Admin.GET("/campaign_trx/users", handler.GetUsers)
	echoGroup.Admin.GET("/campaign_trx/point", cmpHandler.GetUserPoint)
	echoGroup.Admin.GET("/campaign_trx/point/history", cmpHandler.GetUserPointHistory)
}

// GetUsers a handler to create a campaignTrx
func (cmpgnTrx *CampaigntrxsHandler) GetUsers(c echo.Context) error {
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
	requestLogger.Info("Start to get campaignTrx users.")
	data, counter, err := cmpgnTrx.CampaigntrxUseCase.GetUsers(c, payload)

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

	requestLogger.Info("End of get campaignTrx users.")

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
		return http.StatusInternalServerError
	}
}
