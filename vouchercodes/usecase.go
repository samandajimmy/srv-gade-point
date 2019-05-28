package vouchercodes

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the vouchercode's usecases
type UseCase interface {
	GetVoucherCodeHistory(echo.Context, map[string]interface{}) ([]models.VoucherCode, string, error)
	GetVoucherCodes(echo.Context, map[string]interface{}) ([]models.VoucherCode, string, error)
	VoucherCodeRedeem(echo.Context, *models.PayloadValidator) (*models.VoucherCode, error)
	GetVoucherCodeByCriteria(echo.Context, map[string]interface{}) ([]models.VoucherCode, string, error)
}
