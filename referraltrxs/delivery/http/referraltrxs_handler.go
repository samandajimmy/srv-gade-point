package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referraltrxs"
	"gade/srv-gade-point/services"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// ReferralTrxHandler represent the httphandler for referralTrx
type ReferralTrxHandler struct {
	ReferralTrxUseCase referraltrxs.UseCase
}

// NewReferralTrxHandler represent to register referralTrx endpoint
func NewReferralTrxHandler(echoGroup models.EchoGroup, us referraltrxs.UseCase) {
	handler := &ReferralTrxHandler{
		ReferralTrxUseCase: us,
	}

	// End Point For External
	echoGroup.API.GET("/milestone/:cif", handler.getMilestone)
	echoGroup.API.GET("/ranking", handler.getRanking)
}

func (rfr *ReferralTrxHandler) getMilestone(c echo.Context) error {
	respErrors := &models.ResponseErrors{}
	response    = models.Response{}
	logger     := models.RequestLogger{}

	// metric monitoring
	go services.AddMetric("get_milestone")

	requestLogger := logger.GetRequestLogger(c, nil)
	requestLogger.Info("Start to get milestone.")

	responseData, err := rfr.ReferralTrxUseCase.GetMilestone(c, c.Param("cif"))

	if err != nil {
		respErrors.SetTitle(models.MessageDataNotFound)
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End to get milestone for client")
		// metric monitoring
		go services.AddMetric("get_milestone_error")

		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	response.SetResponse(responseData, respErrors)
	requestLogger.Info("End of get milestone.")
	// metric monitoring
	go services.AddMetric("get_milestone_success")

	return c.JSON(getStatusCode(err), response)
}

func (rfr *ReferralTrxHandler) getRanking(c echo.Context) error {
	// metric monitoring
	go services.AddMetric("get_all_vouchers")

	var pl models.RankingPayload
	respErrors := &models.ResponseErrors{}
	response = models.Response{}
	logger := models.RequestLogger{}
	logger.DataLog(c, pl).Info("Start to get all ranking for client")
	err := c.Bind(&pl)

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End to get all ranking for client")

		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if err = c.Validate(pl); err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End to get all ranking for client")

		return c.JSON(http.StatusBadRequest, response)
	}

	responseData, err := rfr.ReferralTrxUseCase.GetRanking(c, pl)

	if err != nil {
		// metric monitoring error
		go services.AddMetric("get_all_ranking_error")

		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		response.Status = models.StatusError
		response.Message = err.Error()

		logger.DataLog(c, response).Info("End to get all ranking for client")
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.SetResponse(responseData, respErrors)
	logger.DataLog(c, response).Info("End to get all ranking for client")

	// metric monitoring success
	go services.AddMetric("get_all_ranking_success")

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
