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
	GetVoucher(ctx context.Context, name string, status string, startDate string, endDate string) (interface{}, error)
}
