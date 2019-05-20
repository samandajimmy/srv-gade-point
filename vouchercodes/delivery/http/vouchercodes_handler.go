package http

import (
	"gade/srv-gade-point/models"
	vouchercodes "gade/srv-gade-point/vouchercodes"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// VoucherCodesHandler represent the httphandler for vouchercodes
type VoucherCodesHandler struct {
	VoucherCodeUseCase vouchercodes.UseCase
}

// NewVoucherCodesHandler represent to view voucher history endpoint
func NewVoucherCodesHandler(echoGroup models.EchoGroup, vcu vouchercodes.UseCase) {
	handler := &VoucherCodesHandler{
		VoucherCodeUseCase: vcu,
	}

	//End Point For CMS
	echoGroup.Admin.GET("/voucher_codes/voucher/history", handler.GetVoucherCodeHistory)

	//End Point For External
	echoGroup.API.GET("/voucher_codes/voucher/history", handler.GetVoucherCodeHistory)

}

// GetVoucherCodeHistory a handler to create a vouchercodes
func (VchrCode *VoucherCodesHandler) GetVoucherCodeHistory(c echo.Context) error {
	response = models.Response{}
	userID := c.QueryParam("userId")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	if pageStr == "" {
		pageStr = "0"
	}

	if limitStr == "" {
		limitStr = "0"
	}

	payload := map[string]interface{}{
		"userId": userID,
		"page":   pageStr,
		"limit":  limitStr,
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
	requestLogger.Info("Start to get voucher history.")
	data, counter, err := VchrCode.VoucherCodeUseCase.GetVoucherCodeHistory(c, payload)

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

	requestLogger.Info("End of get voucher history.")

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
