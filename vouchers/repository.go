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
	GetVouchers(ctx context.Context, name string, status string, startDate string, endDate string, page int32, limit int32) ([]*models.Voucher, error)
	GetVouchersExternal(ctx context.Context, name string, startDate string, endDate string, page int32, limit int32) ([]*models.VoucherDetail, error)
	GetVoucher(ctx context.Context, voucherId string) (*models.Voucher, error)
	GetVoucherExternal(ctx context.Context, voucherId string) (*models.VoucherDetail, error)
	CountVouchers(ctx context.Context, status string, expired bool) (int, error)
	DeleteVoucher(ctx context.Context, id int64) error
	GetVouchersUser(ctx context.Context, userId string, status string, page int32, limit int32) ([]*models.VoucherUser, error)
	CountPromoCode(ctx context.Context, status string, userId string) (int, error)
	UpdatePromoCodeBought(ctx context.Context, voucherId string, userId string) (int64, string, string, error)
}
