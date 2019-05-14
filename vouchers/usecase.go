package vouchers

import (
	"gade/srv-gade-point/models"
	"mime/multipart"

	"github.com/labstack/echo"
)

// UseCase represent the voucher's usecases
type UseCase interface {
	CreateVoucher(echo.Context, *models.Voucher) error
	UpdateVoucher(echo.Context, int64, *models.UpdateVoucher) error
	UploadVoucherImages(echo.Context, *multipart.FileHeader) (string, error)
	GetVouchersAdmin(echo.Context, map[string]interface{}) ([]*models.Voucher, string, error)
	GetVoucherAdmin(echo.Context, string) (*models.Voucher, error)
	GetVouchers(echo.Context, map[string]interface{}) ([]*models.Voucher, string, error)
	GetVoucher(echo.Context, string) (*models.Voucher, error)
	VoucherBuy(echo.Context, *models.PayloadVoucherBuy) (*models.VoucherCode, error)
	GetVouchersUser(echo.Context, map[string]interface{}) ([]models.VoucherCode, string, error)
	VoucherValidate(echo.Context, *models.PayloadValidateVoucher) (*models.ResponseValidateVoucher, error)
	VoucherRedeem(echo.Context, *models.PayloadValidateVoucher) (*models.VoucherCode, error)
	UpdateStatusBasedOnStartDate() error
}
