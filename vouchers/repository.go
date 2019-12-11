package vouchers

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the voucher's repository contract
type Repository interface {
	UpdateVoucherStock(c echo.Context, voucherId string) error
	CreateVoucher(echo.Context, *models.Voucher) error
	CreatePromoCode(echo.Context, []*models.VoucherCode) error
	UpdateVoucher(echo.Context, int64, *models.UpdateVoucher) error
	GetVouchersAdmin(echo.Context, map[string]interface{}) ([]*models.Voucher, error)
	GetVoucherAdmin(echo.Context, string) (*models.Voucher, error)
	GetVouchers(echo.Context) ([]*models.Voucher, error)
	GetVoucher(echo.Context, string) (*models.Voucher, error)
	UpdatePromoCodeBought(echo.Context, string, string) (*models.VoucherCode, error)
	BookVoucherCode(echo.Context, *models.PayloadVoucherBuy) (*models.VoucherCode, error)
	GetVouchersUser(echo.Context, map[string]interface{}) ([]models.VoucherCode, error)
	CountVouchers(echo.Context, bool) (int, error)
	DeleteVoucher(echo.Context, int64) error
	CountPromoCode(echo.Context, map[string]interface{}) (int, error)
	UpdatePromoCodeRedeemed(echo.Context, string, string, string) (*models.VoucherCode, error)
	GetVoucherCode(echo.Context, *models.PayloadValidator) (*models.VoucherCode, string, error)
	UpdateExpiryDate(echo.Context) error
	UpdateStatusBasedOnStartDate() error
	CountBoughtVoucher(echo.Context, string, string) (int64, error)
}
