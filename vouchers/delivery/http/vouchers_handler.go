package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
)

var response = models.Response{} // Response represent the response

// VouchersHandler represent the httphandler for vouchers
type VouchersHandler struct {
	VoucherUseCase vouchers.UseCase
}

// NewVouchersHandler represent to register vouchers endpoint
func NewVouchersHandler(echoGroup models.EchoGroup, us vouchers.UseCase) {
	handler := &VouchersHandler{
		VoucherUseCase: us,
	}

	// End Point For CMS
	echoGroup.Admin.POST("/vouchers", handler.CreateVoucher)
	echoGroup.Admin.GET("/vouchers", handler.GetVouchersAdmin)
	echoGroup.Admin.GET("/vouchers/:id", handler.GetVoucherAdmin)
	echoGroup.Admin.POST("/vouchers/upload", handler.UploadVoucherImages)
	echoGroup.Admin.PUT("/vouchers/status/:id", handler.UpdateStatusVoucher)

	// End Point For External
	echoGroup.API.GET("/vouchers", handler.GetVouchers)
	echoGroup.API.GET("/hidden/vouchers/:id", handler.GetVoucher)
	echoGroup.API.POST("/vouchers/buy", handler.VoucherBuy)
	echoGroup.API.POST("/hidden/vouchers/redeem", handler.VoucherRedeem)
	echoGroup.API.GET("/vouchers/user", handler.GetVouchersUser)
	echoGroup.API.GET("/vouchers/user/history", handler.GetHistoryVouchersUser)
	echoGroup.API.POST("/hidden/vouchers/validate", handler.VoucherValidate)
}

// CreateVoucher Create new voucher and generate promo code by stock
func (vchr *VouchersHandler) CreateVoucher(c echo.Context) error {
	var voucher models.Voucher
	respErrors := &models.ResponseErrors{}
	response = models.Response{}
	logger := models.RequestLogger{}

	err := c.Bind(&voucher)
	logger.DataLog(c, voucher).Info("Start to create a voucher.")

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End of create a voucher")

		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	err = vchr.VoucherUseCase.CreateVoucher(c, &voucher)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End of create a voucher")

		return c.JSON(getStatusCode(err), response)
	}

	response.SetResponse(voucher, respErrors)
	logger.DataLog(c, response).Info("End of create a voucher")

	return c.JSON(getStatusCode(err), response)
}

// UpdateStatusVoucher Update status voucher ACTIVE or INACTIVE
func (vchr *VouchersHandler) UpdateStatusVoucher(c echo.Context) error {
	var updateVoucher *models.UpdateVoucher
	respErrors := &models.ResponseErrors{}
	response = models.Response{}
	logger := models.RequestLogger{}
	err := c.Bind(updateVoucher)
	logger.DataLog(c, *updateVoucher).Info("Start to update voucher.")

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End of update a voucher.")

		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	updateVoucher.VoucherID, err = strconv.ParseInt(c.Param("id"), 10, 64)

	if err != nil {
		respErrors.SetTitle(http.StatusText(http.StatusBadRequest))
		response.SetResponse("", respErrors)
		logrus.Debug(err)
		logger.DataLog(c, response).Info("End of update a voucher.")

		return c.JSON(http.StatusBadRequest, response)
	}

	err = vchr.VoucherUseCase.UpdateVoucher(c, updateVoucher.VoucherID, updateVoucher)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logrus.Debug(err)
		logger.DataLog(c, response).Info("End of update a voucher.")

		return c.JSON(getStatusCode(err), response)
	}

	response.SetResponse("", respErrors)
	logger.DataLog(c, response).Info("End of update a voucher.")

	return c.JSON(getStatusCode(err), response)
}

// UploadVoucherImages Upload image voucher
func (vchr *VouchersHandler) UploadVoucherImages(c echo.Context) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	response = models.Response{}
	file, err := c.FormFile("file")

	requestLogger.Info("Start to upload an voucher image.")
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	path, err := vchr.VoucherUseCase.UploadVoucherImages(c, file)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if path != "" {
		response.Data = models.PathVoucher{ImageURL: path}
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageUploadSuccess
	requestLogger.Info("End of upload an voucher image.")

	return c.JSON(getStatusCode(err), response)
}

// GetVouchersAdmin Get all voucher by param name, status, start date and end date for admin
func (vchr *VouchersHandler) GetVouchersAdmin(c echo.Context) error {
	response = models.Response{}
	name := c.QueryParam("name")
	status := c.QueryParam("status")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	// validate page and limit string input
	if pageStr == "" {
		pageStr = "0"
	}

	if limitStr == "" {
		limitStr = "0"
	}

	// prepare payload for logger
	payload := map[string]interface{}{
		"name":      name,
		"status":    status,
		"page":      pageStr,
		"limit":     limitStr,
		"startDate": startDate,
		"endDate":   endDate,
	}

	logger := models.RequestLogger{
		Payload: payload,
	}

	requestLogger := logger.GetRequestLogger(c, payload)
	requestLogger.Info("Start to get all voucher for admin")

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

	if startDate != "" && !dateFmtRgx.MatchString(startDate) {
		requestLogger.Debug(models.ErrStartDateFormat)
		response.Status = models.StatusError
		response.Message = models.ErrStartDateFormat.Error()

		return c.JSON(http.StatusBadRequest, response)
	}

	if endDate != "" && !dateFmtRgx.MatchString(endDate) {
		requestLogger.Debug(models.ErrEndDateFormat)
		response.Status = models.StatusError
		response.Message = models.ErrEndDateFormat.Error()

		return c.JSON(http.StatusBadRequest, response)
	}

	payload["page"] = page
	payload["limit"] = limit
	responseData, totalCount, err := vchr.VoucherUseCase.GetVouchersAdmin(c, payload)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if len(responseData) > 0 {
		response.Data = responseData
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.TotalCount = totalCount
	requestLogger.Info("End of get all voucher for admin")

	return c.JSON(getStatusCode(err), response)
}

// GetVoucherAdmin Get detail voucher by voucherId for admin
func (vchr *VouchersHandler) GetVoucherAdmin(c echo.Context) error {
	response = models.Response{}
	voucherID := c.Param("id")
	logger := models.RequestLogger{
		Payload: map[string]interface{}{
			"voucherID": voucherID,
		},
	}
	requestLogger := logger.GetRequestLogger(c, nil)
	requestLogger.Info("Start to get voucher detail for admin.")
	responseData, err := vchr.VoucherUseCase.GetVoucherAdmin(c, voucherID)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageDataNotFound
		return c.JSON(getStatusCode(err), response)
	}

	if (&models.Voucher{}) != responseData {
		response.Data = responseData
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	requestLogger.Info("End of get voucher detail for admin.")

	return c.JSON(getStatusCode(err), response)
}

// GetVouchers Get all voucher by param name, status, start date and end date
func (vchr *VouchersHandler) GetVouchers(c echo.Context) error {
	var pl models.PayloadList
	respErrors := &models.ResponseErrors{}
	response = models.Response{}
	logger := models.RequestLogger{}
	logger.DataLog(c, pl).Info("Start to get all voucher for client")
	err := c.Bind(&pl)

	if err != nil {
		respErrors.SetTitle(models.MessageUnprocessableEntity)
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End to get all voucher for client")

		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if err = c.Validate(pl); err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		logger.DataLog(c, response).Info("End to get all voucher for client")

		return c.JSON(http.StatusBadRequest, response)
	}

	responseData, respErrors, err := vchr.VoucherUseCase.GetVouchers(c)

	if err != nil {
		respErrors.SetTitle(err.Error())
		response.SetResponse("", respErrors)
		response.Status = models.StatusError
		response.Message = err.Error()

		logger.DataLog(c, response).Info("End to get all voucher for client")
		return c.JSON(getStatusCode(err), response)
	}

	response.SetResponse(responseData, respErrors)
	logger.DataLog(c, response).Info("End to get all voucher for client")

	return c.JSON(getStatusCode(err), response)
}

// GetVouchers Get all history vouchers by param name, start date and end date
func (vchr *VouchersHandler) GetHistoryVouchersUser(c echo.Context) error {
	response = models.Response{}
	userID := c.QueryParam("userId")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	// validate page and limit string input
	if pageStr == "" {
		pageStr = "0"
	}

	if limitStr == "" {
		limitStr = "0"
	}

	// prepare payload for logger
	payload := map[string]interface{}{
		"userID": userID,
		"page":   pageStr,
		"limit":  limitStr,
	}

	logger := models.RequestLogger{
		Payload: payload,
	}

	requestLogger := logger.GetRequestLogger(c, payload)
	requestLogger.Info("Start to get all voucher for client")

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

	payload["page"] = page
	payload["limit"] = limit
	responseData, err := vchr.VoucherUseCase.GetHistoryVouchersUser(c)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData

	return c.JSON(getStatusCode(err), response)
}

// GetVoucher Get detail voucher by voucherId
func (vchr *VouchersHandler) GetVoucher(c echo.Context) error {
	response = models.Response{}
	voucherID := c.Param("id")
	logger := models.RequestLogger{
		Payload: map[string]interface{}{
			"voucherID": voucherID,
		},
	}

	requestLogger := logger.GetRequestLogger(c, nil)
	requestLogger.Info("Start to get detail voucher for client")
	responseData, err := vchr.VoucherUseCase.GetVoucher(c, voucherID)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageDataNotFound
		return c.JSON(getStatusCode(err), response)
	}

	if (&models.Voucher{}) != responseData {
		response.Data = responseData
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess

	requestLogger.Info("End of get detail voucher for client")
	return c.JSON(getStatusCode(err), response)
}

// GetVouchersUser Get all promo code voucher by user id and status bought
func (vchr *VouchersHandler) GetVouchersUser(c echo.Context) error {
	response = models.Response{}
	userID := c.QueryParam("userId")
	pageStr := c.QueryParam("page")
	limitStr := c.QueryParam("limit")

	// validate page and limit string input
	if pageStr == "" {
		pageStr = "0"
	}

	if limitStr == "" {
		limitStr = "0"
	}

	// prepare payload for logger
	payload := map[string]interface{}{
		"userID": userID,
		"page":   pageStr,
		"limit":  limitStr,
	}

	logger := models.RequestLogger{
		Payload: payload,
	}

	requestLogger := logger.GetRequestLogger(c, payload)
	requestLogger.Info("Start to get all voucher for client")

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

	payload["page"] = page
	payload["limit"] = limit
	responseData, err := vchr.VoucherUseCase.GetVouchersUser(c, payload)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData

	return c.JSON(getStatusCode(err), response)
}

// VoucherBuy is a handler to provide and endpoint to buy voucher with point
func (vchr *VouchersHandler) VoucherBuy(c echo.Context) error {
	var voucher models.PayloadVoucherBuy
	response = models.Response{}

	if err := c.Bind(&voucher); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, voucher)
	requestLogger.Info("Start to buy a voucher")
	responseData, err := vchr.VoucherUseCase.VoucherBuy(c, &voucher)

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
	requestLogger.Info("End of buy a voucher")

	return c.JSON(getStatusCode(err), response)
}

// VoucherValidate is a handler to provide and endpoint to validate voucher before reedem
func (vchr *VouchersHandler) VoucherValidate(c echo.Context) error {
	var payloadValidator models.PayloadValidator
	response = models.Response{}

	if err := c.Bind(&payloadValidator); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, payloadValidator)
	requestLogger.Info("Start to validate a voucher")
	responseData, err := vchr.VoucherUseCase.VoucherValidate(c, &payloadValidator, nil)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	if len(responseData) > 0 {
		response.Data = responseData
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	requestLogger.Info("End of validate a voucher")

	return c.JSON(getStatusCode(err), response)
}

// VoucherRedeem is a handler to provide and endpoint to reedem voucher
func (vchr *VouchersHandler) VoucherRedeem(c echo.Context) error {
	var voucher models.PayloadValidator
	response = models.Response{}

	if err := c.Bind(&voucher); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, voucher)
	requestLogger.Info("Start to redeem a voucher")
	responseData, err := vchr.VoucherUseCase.VoucherRedeem(c, &voucher)

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
	requestLogger.Info("End of redeem a voucher")

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
