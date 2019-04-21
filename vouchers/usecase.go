package vouchers

import (
	"context"
	"gade/srv-gade-point/models"
	"mime/multipart"

	"github.com/labstack/echo"
)

// UseCase represent the voucher's usecases
type UseCase interface {
	CreateVoucher(context.Context, *models.Voucher) error
	UpdateVoucher(ctx context.Context, id int64, updateVoucher *models.UpdateVoucher) error
	UploadVoucherImages(*multipart.FileHeader) (string, error)
	GetVouchersAdmin(ctx context.Context, name string, status string, startDate string, endDate string, page int, limit int) ([]*models.Voucher, string, error)
	GetVoucherAdmin(ctx context.Context, voucherID string) (*models.Voucher, error)
	GetVouchers(ctx context.Context, name string, status string, startDate string, endDate string, page int, limit int) ([]*models.Voucher, string, error)
	GetVoucher(ctx context.Context, voucherID string) (*models.Voucher, error)
	VoucherBuy(ctx context.Context, ech echo.Context, m *models.PayloadVoucherBuy) (*models.PromoCode, error)
	GetVouchersUser(ctx context.Context, userID string, status string, page int, limit int) ([]models.PromoCode, string, error)
	VoucherValidate(ctx context.Context, m *models.PayloadValidateVoucher) (*models.ResponseValidateVoucher, error)
	VoucherRedeem(ctx context.Context, m *models.PayloadValidateVoucher) (*models.PromoCode, error)
	UpdateStatusBasedOnStartDate() error
}
