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

// Response represent the response
var response = models.Response{}

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
	e.POST("/vouchers/uploadImage", handler.UploadVoucherImages)
	e.GET("/vouchers", handler.GetVouchers)
}

func (a *VouchersHandler) CreateVoucher(c echo.Context) error {
	var voucher models.Voucher
	err := c.Bind(&voucher)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if ok, err := isRequestValid(&voucher); !ok {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusBadRequest, response)
	}

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.VoucherUseCase.CreateVoucher(ctx, &voucher)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MassageSaveSuccess
	response.Data = voucher
	return c.JSON(http.StatusCreated, response)
}

func (a *VouchersHandler) UpdateStatusVoucher(c echo.Context) error {

	updateVoucher := new(models.UpdateVoucher)

	if err := c.Bind(updateVoucher); err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if ok, err := isRequestValid(updateVoucher); !ok {
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

	err := a.VoucherUseCase.UpdateVoucher(ctx, int64(id), updateVoucher)

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

func (a *VouchersHandler) UploadVoucherImages(c echo.Context) error {

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	file, err := c.FormFile("file")
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(getStatusCode(err), response)
	}

	path, err := a.VoucherUseCase.UploadVoucherImages(file)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MassageUploadSuccess
	response.Data = models.PathVoucher{ImageUrl: path}
	return c.JSON(http.StatusOK, response)
}

func (a *VouchersHandler) GetVouchers(c echo.Context) error {

	name := c.QueryParam("name")
	status := c.QueryParam("status")
	startDate := c.QueryParam("startDate")
	endDate := c.QueryParam("endDate")

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, err := a.VoucherUseCase.GetVoucher(ctx, name, status, startDate, endDate)

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
