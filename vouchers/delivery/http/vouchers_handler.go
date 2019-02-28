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
	response = models.Response{} // Response represent the response
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
	e.GET("/vouchers", handler.GetVouchers)
	e.GET("/vouchers/monitor", handler.GetVouchersMonitoring)
}

//Create new voucher and generate promo code by stock
func (a *VouchersHandler) CreateVoucher(c echo.Context) error {
	var voucher models.Voucher
	err := c.Bind(&voucher)
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		response.TotalCount = ""
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	if ok, err := isRequestValid(&voucher); !ok {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		response.TotalCount = ""
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
		response.TotalCount = ""
		return c.JSON(http.StatusUnprocessableEntity, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MassageSaveSuccess
	response.Data = voucher
	response.TotalCount = ""
	return c.JSON(http.StatusCreated, response)
}

//Update status voucher ACTIVE or INACTIVE
func (a *VouchersHandler) UpdateStatusVoucher(c echo.Context) error {

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

	err := a.VoucherUseCase.UpdateVoucher(ctx, int64(id), updateVoucher)

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

//Upload image voucher
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
		response.TotalCount = ""
		return c.JSON(getStatusCode(err), response)
	}

	path, err := a.VoucherUseCase.UploadVoucherImages(file)
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

//Get all voucher by param name, status, start date and end date
func (a *VouchersHandler) GetVouchers(c echo.Context) error {

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

	res, totalCount, err := a.VoucherUseCase.GetVouchers(ctx, name, status, startDate, endDate, int32(page), int32(limit))

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		response.TotalCount = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.StatusSuccess
	response.Data = res
	response.TotalCount = totalCount
	return c.JSON(http.StatusOK, response)
}

//Get monitoring voucher stock amount, stock avaliable, stock bought,stock redeemed, expired
func (a *VouchersHandler) GetVouchersMonitoring(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))

	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	res, totalCount, err := a.VoucherUseCase.GetVouchersMonitoring(ctx, int32(page), int32(limit))
	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()
		response.Data = ""
		response.TotalCount = ""
		return c.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.StatusSuccess
	response.Data = res
	response.TotalCount = totalCount
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
