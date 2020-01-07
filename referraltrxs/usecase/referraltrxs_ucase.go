package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referraltrxs"
	"github.com/labstack/echo"
	"os"
	"strconv"
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

func (rfr *referralTrxUseCase) GetMilestone(c echo.Context, pl models.MilestonePayload) (*models.Milestone, error) {
	logger 		   := models.RequestLogger{}
	requestLogger  := logger.GetRequestLogger(c, nil)
	milestone, err := rfr.referralTrxRepo.GetMilestone(c, pl.CIF)

	if err != nil {
		requestLogger.Debug(models.ErrMilestone)

		return nil, models.ErrMilestone
	}

	ranking, err := rfr.referralTrxRepo.GetRankingByReferralCode(c, pl.ReferralCode)

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

	return milestone, nil
}

func (rfr *referralTrxUseCase) GetRanking(c echo.Context, payload models.RankingPayload) ([]models.Ranking, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// get data
	data, err :=  rfr.referralTrxRepo.GetRanking(c, payload.ReferralCode)

	if err != nil {
		requestLogger.Debug(models.ErrRanking)

		return nil, models.ErrRanking
	}

	return data, err
}
