package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchercodes"

	"github.com/labstack/echo"
)

type voucherCodeUseCase struct {
	voucherCodeRepo vouchercodes.Repository
}

// NewVoucherCodeUseCase will create new an NewVoucherCodeUseCase object representation of vouchercode.UseCase interface
func NewVoucherCodeUseCase(vchrCodesRepo vouchercodes.Repository) vouchercodes.UseCase {
	return &voucherCodeUseCase{
		voucherCodeRepo: vchrCodesRepo,
	}
}

func (vchrCodeUs *voucherCodeUseCase) GetVoucherCodeHistory(c echo.Context, payload map[string]interface{}) ([]models.VoucherCode, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// count voucher by userId
	counter, err := vchrCodeUs.voucherCodeRepo.CountVoucherCode(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrUsersNA)

		return nil, "", models.ErrUsersNA
	}

	// get voucher code history data
	data, err := vchrCodeUs.voucherCodeRepo.GetVoucherCodeHistory(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherHistory)

		return nil, "", models.ErrGetVoucherHistory
	}

	return data, counter, err
}

func (vchrCodeUs *voucherCodeUseCase) GetVoucherCodes(c echo.Context, payload map[string]interface{}) ([]models.VoucherCode, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// count voucher by voucherId
	counter, err := vchrCodeUs.voucherCodeRepo.CountVoucherCodeByVoucherID(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrUsersNA)

		return nil, "", models.ErrUsersNA
	}

	// get voucher codes
	data, err := vchrCodeUs.voucherCodeRepo.GetVoucherCodes(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetVoucherCodes)

		return nil, "", models.ErrGetVoucherCodes
	}

	return data, counter, err
}
