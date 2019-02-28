package vouchers

import (
	"context"
	"gade/srv-gade-point/models"
)

// Repository represent the voucher's repository contract
type Repository interface {
	CreateVoucher(ctx context.Context, a *models.Voucher) error
	UpdateVoucher(ctx context.Context, id int64, updateVoucher *models.UpdateVoucher) error
	GetVouchers(ctx context.Context, name string, status string, startDate string, endDate string, page int32, limit int32) ([]*models.Voucher, error)
	CreatePromoCode(ctx context.Context, promoCode []*models.PromoCode) error
	GetVouchersMonitoring(ctx context.Context, page int32, limit int32) ([]*models.VouchersMonitoring, error)
	CountVouchers(ctx context.Context, status string) (int, error)
	DeleteVoucher(ctx context.Context, id int64) error
}
