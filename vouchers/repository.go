package vouchers

import (
	"context"
	"gade/srv-gade-point/models"
)

// Repository represent the voucher's repository contract
type Repository interface {
	CreateVoucher(ctx context.Context, a *models.Voucher) error
	CreatePromoCode(ctx context.Context, promoCode []*models.PromoCode) error
	UpdateVoucher(ctx context.Context, id int64, updateVoucher *models.UpdateVoucher) error
	GetVouchersAdmin(ctx context.Context, name string, status string, startDate string, endDate string, page int, limit int) ([]*models.Voucher, error)
	GetVoucherAdmin(ctx context.Context, voucherID string) (*models.Voucher, error)
	GetVouchers(ctx context.Context, name string, startDate string, endDate string, page int, limit int) ([]*models.Voucher, error)
	GetVoucher(ctx context.Context, voucherID string) (*models.Voucher, error)
	UpdatePromoCodeBought(ctx context.Context, voucherID string, userID string) (*models.PromoCode, error)
	GetVouchersUser(ctx context.Context, userID string, status string, page int, limit int) ([]models.PromoCode, error)
	CountVouchers(ctx context.Context, name string, status string, startDate string, endDate string, expired bool) (int, error)
	DeleteVoucher(ctx context.Context, id int64) error
	CountPromoCode(ctx context.Context, status string, userID string) (int, error)
	VoucherCheckExpired(ctx context.Context, voucherID string) error
	UpdatePromoCodeRedeemed(ctx context.Context, voucherID string, userID string) (*models.PromoCode, error)
	GetVoucherCode(ctx context.Context, voucherCode string, userID string) (*models.PromoCode, string, error)
	UpdateExpiryDate(ctx context.Context) error
	UpdateStatusBasedOnStartDate() error
}
