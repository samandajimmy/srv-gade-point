package vouchercodes

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the vouchercode's repository contract
type Repository interface {
	GetVoucherCodeHistory(echo.Context, map[string]interface{}) ([]models.VoucherCode, error)
	CountVoucherCode(echo.Context, map[string]interface{}) (string, error)
	GetVoucherCodes(echo.Context, map[string]interface{}) ([]models.VoucherCode, error)
	CountVoucherCodeByVoucherID(echo.Context, map[string]interface{}) (string, error)
}
