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
	echoGroup.Admin.GET("/voucher-codes/:voucherId", handler.GetVoucherCodes)
	echoGroup.Admin.POST("/voucher-codes/import", handler.ImportVoucherCodes)
	echoGroup.Admin.GET("/voucher-codes/voucher/history", handler.GetVoucherCodeHistory)
	echoGroup.Admin.GET("/voucher-codes/:id", handler.GetVoucherCodes)
	echoGroup.Admin.GET("/voucher-codes/bought", handler.GetBoughtVoucherCode)
	echoGroup.Admin.POST("/voucher-codes/redeem", handler.VoucherCodeRedeem)

	//End Point For External
	echoGroup.API.GET("/hidden/voucher-codes/voucher/history", handler.GetVoucherCodeHistory)

}

// GetVoucherCodeHistory a handler to get vouchercodes history
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

// GetVoucherCodes a handler to get vouchercodes
func (VchrCode *VoucherCodesHandler) GetVoucherCodes(c echo.Context) error {
	response = models.Response{}
	voucherID := c.Param("voucherId")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	if pageStr == "" {
		pageStr = "0"
	}

	if limitStr == "" {
		limitStr = "0"
	}

	payload := map[string]interface{}{
		"voucherId": voucherID,
		"page":      pageStr,
		"limit":     limitStr,
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
	requestLogger.Info("Start to get voucher codes.")
	data, counter, err := VchrCode.VoucherCodeUseCase.GetVoucherCodes(c, payload)

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

	requestLogger.Info("End of get voucher codes.")

	return c.JSON(getStatusCode(err), response)
}

// VoucherCodeRedeem is a handler to provide and endpoint to reedem voucher code
func (VchrCode *VoucherCodesHandler) VoucherCodeRedeem(c echo.Context) error {
	var voucher models.PayloadValidator
	response = models.Response{}

	if err := c.Bind(&voucher); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, voucher)
	requestLogger.Info("Start to redeem a voucher code")
	responseData, err := VchrCode.VoucherCodeUseCase.VoucherCodeRedeem(c, &voucher)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if (&models.VoucherCode{}) != responseData {
		response.Data = responseData
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	requestLogger.Info("End of redeem a voucher code")

	return c.JSON(getStatusCode(err), response)
}

// ImportVoucherCodes a handler to import voucher json file
func (VchrCode *VoucherCodesHandler) ImportVoucherCodes(c echo.Context) error {
	response = models.Response{}
	file, _ := c.FormFile("file")
	voucherID := c.FormValue("voucherId")
	payload := map[string]interface{}{"voucherId": voucherID}
	logger := models.RequestLogger{Payload: payload}
	requestLogger := logger.GetRequestLogger(c, nil)
	requestLogger.Info("Start to import voucher codes.")
	_, err := VchrCode.VoucherCodeUseCase.ImportVoucherCodes(c, file, voucherID)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageUploadSuccess
	requestLogger.Info("End of import voucher codes.")

	return c.JSON(getStatusCode(err), response)

}

// GetBoughtVoucherCode is a handler to provide and endpoint to search voucher code
func (VchrCode *VoucherCodesHandler) GetBoughtVoucherCode(c echo.Context) error {
	response = models.Response{}
	promoCode := c.QueryParam("promoCode")
	userID := c.QueryParam("userId")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	payload := map[string]interface{}{
		"promoCode": promoCode,
		"userId":    userID,
		"page":      pageStr,
		"limit":     limitStr,
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
	requestLogger.Info("Start to get voucher codes.")
	data, counter, err := VchrCode.VoucherCodeUseCase.GetBoughtVoucherCode(c, payload)

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

	requestLogger.Info("End of get voucher codes.")
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
