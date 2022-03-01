package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referraltrxs"
	"os"
	"strconv"

	"github.com/labstack/echo"
)

type referralTrxUseCase struct {
	referralTrxRepo referraltrxs.Repository
}

// NewReferralTrxUseCase will create new an referralTrxUseCase object representation of referraltrxs.UseCase interface
func NewReferralTrxUseCase(referralTrxRepo referraltrxs.Repository) referraltrxs.UseCase {
	return &referralTrxUseCase{
		referralTrxRepo: referralTrxRepo,
	}
}

func (rfr *referralTrxUseCase) UGetMilestone(c echo.Context, pl models.MilestonePayload) (*models.Milestone, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	milestone, err := rfr.referralTrxRepo.RGetMilestone(c, pl)

	if err != nil {
		requestLogger.Debug(models.ErrMilestone)

		return nil, models.ErrMilestone
	}

	ranking, err := rfr.referralTrxRepo.RGetRankingByReferralCode(c, pl.ReferralCode)

	if err != nil {
		requestLogger.Debug(models.ErrMilestone)

		return nil, models.ErrMilestone
	}

	stage := os.Getenv(`STAGE`)
	limitRewardCounter := os.Getenv(`LIMIT_REWARD_COUNTER`)
	milestone.Stages, _ = strconv.ParseInt(stage, 10, 64)
	milestone.LimitRewardCounter, _ = strconv.ParseInt(limitRewardCounter, 10, 64)
	milestone.LimitReward = milestone.Stages * milestone.LimitRewardCounter
	milestone.Ranking = *ranking

	if milestone.TotalRewardCounter > milestone.LimitRewardCounter {
		milestone.TotalRewardCounter = milestone.LimitRewardCounter
		milestone.TotalReward = milestone.LimitRewardCounter * milestone.LimitReward
	}

	return milestone, nil
}

func (rfr *referralTrxUseCase) UGetRanking(c echo.Context, payload models.RankingPayload) ([]*models.Ranking, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// get data
	data, err := rfr.referralTrxRepo.RGetRanking(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrRanking)

		return nil, models.ErrRanking
	}

	return data, err
}
