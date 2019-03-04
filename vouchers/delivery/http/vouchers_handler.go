package http

import (
	"context"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
func NewVouchersHandler(e *echo.Echo, us vouchers.UseCase) {
	handler := &VouchersHandler{
		VoucherUseCase: us,
	}

	e.POST("/vouchers", handler.CreateVoucher)
	e.PUT("/vouchers/status/:id", handler.UpdateStatusVoucher)
	e.POST("/vouchers/upload", handler.UploadVoucherImages)
	e.GET("/voucher", handler.GetVoucher)
	e.GET("/vouchers", handler.GetVouchers)
	e.GET("/vouchers/user", handler.GetVouchersUser)
	e.POST("/voucher/buy", handler.CreateVoucherBuy)

}

// Create new voucher and generate promo code by stock
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

	if ok, err := isRequestValid(&voucher); !ok {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
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
	response.Message = models.MassageSaveSuccess
	response.Data = voucher
	return c.JSON(http.StatusCreated, response)
}

// Update status voucher ACTIVE or INACTIVE
func (vchr *VouchersHandler) UpdateStatusVoucher(c echo.Context) error {

	updateVoucher := new(models.UpdateVoucher)

	if err := c.Bind(updateVoucher); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		response.TotalCount = ""
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if ok, err := isRequestValid(updateVoucher); !ok {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		response.TotalCount = ""
		return c.JSON(http.StatusBadRequest, response)
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
		response.Data = ""
		response.TotalCount = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MassageUpdateSuccess
	response.Data = ""
	response.TotalCount = ""
	return c.JSON(http.StatusOK, response)
}

// Upload image voucher
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
	response.Message = models.MassageUploadSuccess
	response.Data = models.PathVoucher{ImageUrl: path}
	response.TotalCount = ""
	return c.JSON(http.StatusOK, response)
}

// Get all voucher by param name, status, start date and end date
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
	source := c.QueryParam("source")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	sources := viper.GetStringSlice("source.whitelisted.external")

	if checkAccess(source, sources) {
		responseData, totalCount, err = vchr.VoucherUseCase.GetVouchers(ctx, name, status, startDate, endDate, int32(page), int32(limit), source)
		if err != nil {
			response.Status = models.StatusError
			response.Message = err.Error()
			return c.JSON(getStatusCode(err), response)
		}
	} else {
		response.Status = models.StatusError
		response.Message = models.MassageForbiddenError
		return c.JSON(http.StatusForbidden, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.StatusSuccess
	response.Data = responseData
	response.TotalCount = totalCount
	return c.JSON(http.StatusOK, response)
}

// Get detail voucher by voucherId
func (vchr *VouchersHandler) GetVoucher(c echo.Context) error {
	totalCount := ""
	response.Data = ""
	response.TotalCount = ""

	voucherId := c.QueryParam("voucherId")
	source := c.QueryParam("source")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	sources := viper.GetStringSlice("source.whitelisted.external")

	if checkAccess(source, sources) {
		responseData, err = vchr.VoucherUseCase.GetVoucher(ctx, voucherId, source)
		if err != nil {
			response.Status = models.StatusError
			response.Message = models.MessageDataNotFound
			return c.JSON(getStatusCode(err), response)
		}
	} else {
		response.Status = models.StatusError
		response.Message = models.MassageForbiddenError
		return c.JSON(http.StatusForbidden, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.StatusSuccess
	response.Data = responseData
	response.TotalCount = totalCount
	return c.JSON(http.StatusOK, response)
}

// Get all promo code voucher by user id and status bought
func (vchr *VouchersHandler) GetVouchersUser(c echo.Context) error {
	totalCount := ""
	response.Data = ""
	response.TotalCount = ""

	userId := c.QueryParam("userId")
	status := c.QueryParam("status")
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	source := c.QueryParam("source")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	sources := viper.GetStringSlice("source.whitelisted.external")

	if checkAccess(source, sources) {
		responseData, totalCount, err = vchr.VoucherUseCase.GetVouchersUser(ctx, userId, status, int32(page), int32(limit), source)
		if err != nil {
			response.Status = models.StatusError
			response.Message = err.Error()
			return c.JSON(getStatusCode(err), response)
		}
	} else {
		response.Status = models.StatusError
		response.Message = models.MassageForbiddenError
		return c.JSON(http.StatusForbidden, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MassagePointSuccess
	response.Data = responseData
	response.TotalCount = totalCount
	return c.JSON(http.StatusOK, response)
}

// Buy voucher
func (vchr *VouchersHandler) CreateVoucherBuy(c echo.Context) error {
	var voucher models.PayloadVoucherBuy
	response.Data = ""
	response.TotalCount = ""
	err = c.Bind(&voucher)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if ok, err := isRequestValid(&voucher); !ok {
		response.Status = models.StatusError
		response.Message = err.Error()
		return c.JSON(http.StatusBadRequest, response)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	sources := viper.GetStringSlice("source.whitelisted.external")

	if checkAccess(voucher.Source, sources) {
		responseData, err = vchr.VoucherUseCase.CreateVoucherBuy(ctx, &voucher)
		if err != nil {
			response.Status = models.StatusError
			response.Message = err.Error()
			return c.JSON(getStatusCode(err), response)
		}
	} else {
		response.Status = models.StatusError
		response.Message = models.MassageForbiddenError
		return c.JSON(http.StatusForbidden, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MassagePointSuccess
	response.Data = responseData
	return c.JSON(http.StatusCreated, response)
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

func checkAccess(str string, list []string) bool {
	for _, v := range list {
		if v == str {
			return true
		}
	}
	return false
}
