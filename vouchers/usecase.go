package vouchers

import (
	"gade/srv-gade-point/models"
	"mime/multipart"

	"github.com/labstack/echo"
)

// UseCase represent the voucher's usecases
type VUsecase interface {
	CreateVoucher(echo.Context, *models.Voucher) error
	UpdateVoucher(echo.Context, int64, *models.UpdateVoucher) error
	UploadVoucherImages(echo.Context, *multipart.FileHeader) (string, error)
	GetVouchersAdmin(echo.Context, map[string]interface{}) ([]*models.Voucher, string, error)
	GetVoucherAdmin(echo.Context, string) (*models.Voucher, error)
	GetVouchers(echo.Context) (models.ListVoucher, *models.ResponseErrors, error)
	GetVoucher(echo.Context, string) (*models.Voucher, error)
	VoucherBuy(echo.Context, *models.PayloadVoucherBuy) (*models.VoucherCode, error)
	VoucherGive(echo.Context, *models.PayloadVoucherBuy) (*models.VoucherCode, error)
	GetVouchersUser(echo.Context, map[string]interface{}) (models.ListVoucher, error)
	VoucherValidate(echo.Context, *models.PayloadValidator, *models.VoucherCode) ([]models.Reward, error)
	VoucherRedeem(echo.Context, *models.PayloadValidator) (*models.VoucherCode, error)
	GetVoucherCode(echo.Context, *models.PayloadValidator, bool) (*models.VoucherCode, string, error)
	UpdateStatusBasedOnStartDate() error
	GetHistoryVouchersUser(echo.Context) (models.ListVoucher, error)
}
