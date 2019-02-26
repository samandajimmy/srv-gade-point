package vouchers

import (
	"context"
	"gade/srv-gade-point/models"
)

// Repository represent the voucher's repository contract
type Repository interface {
	CreateVoucher(ctx context.Context, a *models.Voucher) error
	UpdateVoucher(ctx context.Context, id int64, updateVoucher *models.UpdateVoucher) error
	GetVoucher(ctx context.Context, name string, status string, startDate string, endDate string) ([]*models.Voucher, error)
	CreatePromoCode(ctx context.Context, promoCode []*models.PromoCode) error
}
