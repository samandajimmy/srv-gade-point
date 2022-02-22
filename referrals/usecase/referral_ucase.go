package usecase

import (
	"gade/srv-gade-point/referrals"
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
