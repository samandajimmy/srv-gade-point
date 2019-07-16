package usecase

import (
	"gade/srv-gade-point/rewardtrxs"
)

type rewardTrxUseCase struct {
	rewardTrxRepo rewardtrxs.Repository
}

// NewRewardtrxUseCase will create new an rewardtrxUseCase object representation of rewardtrxs.UseCase interface
func NewRewardtrxUseCase(rwdTrxRepo rewardtrxs.Repository) rewardtrxs.UseCase {
	return &rewardTrxUseCase{
		rewardTrxRepo: rwdTrxRepo,
	}
}
