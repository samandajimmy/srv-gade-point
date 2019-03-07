package http

import (
	"context"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

var (
	response     = models.Response{} // Response represent the response
	responseData interface{}
	err          error
)

// VouchersHandler represent the httphandler for vouchers
type VouchersHandler struct {
	VoucherUseCase vouchers.UseCase
}

// NewVouchersHandler represent to register vouchers endpoint
func NewVouchersHandler(echoGroup models.EchoGroup, us vouchers.UseCase) {
	handler := &VouchersHandler{
		VoucherUseCase: us,
	}

	//End Point For CMS
	echoGroup.Admin.POST("/vouchers", handler.CreateVoucher)
	echoGroup.Admin.PUT("/vouchers/status/:id", handler.UpdateStatusVoucher)
	echoGroup.Admin.POST("/vouchers/upload", handler.UploadVoucherImages)
	echoGroup.Admin.GET("/vouchers", handler.GetVouchersAdmin)
	echoGroup.Admin.GET("/voucher", handler.GetVoucherAdmin)
	//End Point For External
	echoGroup.API.GET("/vouchers", handler.GetVouchers)
	echoGroup.API.GET("/voucher", handler.GetVoucher)
	echoGroup.API.GET("/vouchers/user", handler.GetVouchersUser)
	echoGroup.API.POST("/voucher/buy", handler.VoucherBuy)
	echoGroup.API.POST("/voucher/validate", handler.VoucherValidate)
	echoGroup.API.POST("/voucher/redeem", handler.VoucherRedeem)

}

// CreateVoucher Create new voucher and generate promo code by stock
func (vchr *VouchersHandler) CreateVoucher(c echo.Context) error {
	var voucher models.Voucher
	response.Data = ""
	response.TotalCount = ""
	err = c.Bind(&voucher)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = vchr.VoucherUseCase.CreateVoucher(ctx, &voucher)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageSaveSuccess
	response.Data = voucher
	return c.JSON(http.StatusCreated, response)
}

// UpdateStatusVoucher Update status voucher ACTIVE or INACTIVE
func (vchr *VouchersHandler) UpdateStatusVoucher(c echo.Context) error {
	response.Data = ""
	response.TotalCount = ""
	updateVoucher := new(models.UpdateVoucher)

	if err := c.Bind(updateVoucher); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	id, _ := strconv.Atoi(c.Param("id"))

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err := vchr.VoucherUseCase.UpdateVoucher(ctx, int64(id), updateVoucher)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageUpdateSuccess
	return c.JSON(http.StatusOK, response)
}

// UploadVoucherImages Upload image voucher
func (vchr *VouchersHandler) UploadVoucherImages(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		response.TotalCount = ""
		return c.JSON(getStatusCode(err), response)
	}

	path, err := vchr.VoucherUseCase.UploadVoucherImages(file)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		response.TotalCount = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageUploadSuccess
	response.Data = models.PathVoucher{ImageURL: path}
	response.TotalCount = ""
	return c.JSON(http.StatusOK, response)
}

// GetVouchersAdmin Get all voucher by param name, status, start date and end date for admin
func (vchr *VouchersHandler) GetVouchersAdmin(c echo.Context) error {
	totalCount := ""
	response.Data = ""
	response.TotalCount = ""

	name := c.QueryParam("name")
	status := c.QueryParam("status")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	responseData, totalCount, err = vchr.VoucherUseCase.GetVouchersAdmin(ctx, name, status, startDate, endDate, int32(page), int32(limit))
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData
	response.TotalCount = totalCount
	return c.JSON(http.StatusOK, response)
}

// GetVoucherAdmin Get detail voucher by voucherId for admin
func (vchr *VouchersHandler) GetVoucherAdmin(c echo.Context) error {
	response.Data = ""
	response.TotalCount = ""

	voucherID := c.QueryParam("voucherId")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	responseData, err = vchr.VoucherUseCase.GetVoucherAdmin(ctx, voucherID)
	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageDataNotFound
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

// GetVouchers Get all voucher by param name, status, start date and end date
func (vchr *VouchersHandler) GetVouchers(c echo.Context) error {
	totalCount := ""
	response.Data = ""
	response.TotalCount = ""

	name := c.QueryParam("name")
	status := c.QueryParam("status")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	responseData, totalCount, err = vchr.VoucherUseCase.GetVouchers(ctx, name, status, startDate, endDate, int32(page), int32(limit))
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData
	response.TotalCount = totalCount
	return c.JSON(http.StatusOK, response)
}

// GetVoucher Get detail voucher by voucherId
func (vchr *VouchersHandler) GetVoucher(c echo.Context) error {
	response.Data = ""
	response.TotalCount = ""

	voucherID := c.QueryParam("voucherId")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	responseData, err = vchr.VoucherUseCase.GetVoucher(ctx, voucherID)
	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageDataNotFound
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

// Get all promo code voucher by user id and status bought
func (vchr *VouchersHandler) GetVouchersUser(c echo.Context) error {
	response.Data = ""
	response.TotalCount = ""

	userId := c.QueryParam("userId")
	status := c.QueryParam("status")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	responseData, totalCount, err := vchr.VoucherUseCase.GetVouchersUser(ctx, userId, status, int32(page), int32(limit))
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData
	response.TotalCount = totalCount
	return c.JSON(http.StatusOK, response)
}

// VoucherBuy is a handler to provide and endpoint to buy voucher with point
func (vchr *VouchersHandler) VoucherBuy(c echo.Context) error {
	var voucher models.PayloadVoucherBuy
	response.Data = ""
	response.TotalCount = ""

	err = c.Bind(&voucher)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	responseData, err = vchr.VoucherUseCase.VoucherBuy(ctx, &voucher)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData
	return c.JSON(http.StatusCreated, response)
}

// VoucherValidate is a handler to provide and endpoint to validate voucher before reedem
func (a *VouchersHandler) VoucherValidate(c echo.Context) error {
	var voucher models.PayloadValidateVoucher
	response.Data = ""
	response.TotalCount = ""

	err = c.Bind(&voucher)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	responseData, err = a.VoucherUseCase.VoucherValidate(ctx, &voucher)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData
	return c.JSON(http.StatusOK, response)
}

// VoucherRedeem is a handler to provide and endpoint to reedem voucher
func (a *VouchersHandler) VoucherRedeem(c echo.Context) error {
	var voucher models.PayloadValidateVoucher
	response.Data = ""
	response.TotalCount = ""

	err = c.Bind(&voucher)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	responseData, err = a.VoucherUseCase.VoucherRedeem(ctx, &voucher)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessagePointSuccess
	response.Data = responseData
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
