package usecase

import (
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"

	"github.com/labstack/echo"
)

type referralUseCase struct {
	referralRepo referrals.Repository
}

// NewReferralUseCase will create new an rewardtrxUseCase object representation of referrals.UseCase interface
func NewReferralUseCase(rwdTrxRepo referrals.Repository) referrals.UseCase {
	return &referralUseCase{
		referralRepo: rwdTrxRepo,
	}
}

func (reff *referralUseCase) Schedule(c echo.Context, coreTrx models.CoreTrxPayload) ([]models.CoreTrxResponse, error) {
	var responseData []models.CoreTrxResponse
	err := reff.referralRepo.CreateCoreTrx(c, coreTrx)

	if err != nil {
		logger.Make(c, nil).Debug(models.ErrCoreTrxFailed)
		return nil, models.ErrCoreTrxFailed
	}

	return responseData, nil
}
