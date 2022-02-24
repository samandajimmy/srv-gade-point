package usecase

import (
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

func (ref *referralUseCase) PostCoreTrx(c echo.Context, coreTrx []models.CoreTrxPayload) ([]models.CoreTrxResponse, error) {
	var responseData []models.CoreTrxResponse
	err := ref.referralRepo.PostCoreTrx(c, coreTrx)

	if err != nil {
		return nil, models.ErrCoreTrxFailed
	}

	return responseData, nil
}
