package vouchers

import (
	"context"
	"gade/srv-gade-point/models"
	"mime/multipart"
)

// UseCase represent the voucher's usecases
type UseCase interface {
	CreateVoucher(context.Context, *models.Voucher) error
	UpdateVoucher(ctx context.Context, id int64, updateVoucher *models.UpdateVoucher) error
	UploadVoucherImages(*multipart.FileHeader) (string, error)
	GetVouchersAdmin(ctx context.Context, name string, status string, startDate string, endDate string, page int32, limit int32) ([]*models.Voucher, string, error)
	GetVoucherAdmin(ctx context.Context, voucherId string) (*models.Voucher, error)
	GetVouchers(ctx context.Context, name string, status string, startDate string, endDate string, page int32, limit int32) ([]*models.Voucher, string, error)
	GetVoucher(ctx context.Context, voucherId string) (*models.Voucher, error)
	VoucherBuy(ctx context.Context, m *models.PayloadVoucherBuy) (*models.PromoCode, error)
	GetVouchersUser(ctx context.Context, userId string, status string, page int32, limit int32) ([]models.PromoCode, string, error)
	VoucherValidate(ctx context.Context, m *models.PayloadValidateVoucher) (*models.Voucher, error)
	VoucherRedeem(ctx context.Context, m *models.PayloadValidateVoucher) (*models.PromoCode, error)
}
