package vouchercodes

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the vouchercode's repository contract
type VcRepository interface {
	GetVoucherCodeHistory(echo.Context, map[string]interface{}) ([]models.VoucherCode, error)
	CountVoucherCode(echo.Context, map[string]interface{}) (string, error)
	GetVoucherCodes(echo.Context, map[string]interface{}) ([]models.VoucherCode, error)
	CountVoucherCodeByVoucherID(echo.Context, map[string]interface{}) (string, error)
	UpdateVoucherCodeRedeemed(echo.Context, string, string, string) (*models.VoucherCode, error)
	GetBoughtVoucherCode(echo.Context, map[string]interface{}) ([]models.VoucherCode, error)
	CountBoughtVoucherCode(echo.Context, map[string]interface{}) (string, error)
	UpdateVoucherCodeRefID(echo.Context, *models.VoucherCode, string) error
	UpdateVoucherCodeRejected(echo.Context, string) error
	UpdateVoucherCodeSucceeded(echo.Context, *models.RewardPayment) error
	UpdateVoucherCodeInquired(echo.Context, models.VoucherCode, models.PayloadValidator) error
	GetVoucherCodeRefID(echo.Context, string) (*models.VoucherCode, error)
	ValidateVoucherGive(echo.Context, *models.PayloadVoucherBuy) (*models.VoucherCode, error)
}
