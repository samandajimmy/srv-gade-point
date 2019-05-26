package vouchercodes

import (
	"gade/srv-gade-point/models"
	"mime/multipart"

	"github.com/labstack/echo"
)

// UseCase represent the vouchercode's usecases
type UseCase interface {
	GetVoucherCodeHistory(echo.Context, map[string]interface{}) ([]models.VoucherCode, string, error)
	GetVoucherCodes(echo.Context, map[string]interface{}) ([]models.VoucherCode, string, error)
	ImportVoucherCodes(echo.Context, *multipart.FileHeader, string) (string, error)
}
