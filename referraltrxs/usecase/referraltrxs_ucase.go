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

func (rfr *referralTrxUseCase) GetMilestone(c echo.Context, CIF string) (*models.Milestone, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	cif, err := strconv.Atoi(CIF)

	if err != nil {
		requestLogger.Debug(err)

		return nil, models.ErrCIF
	}

	milestone, err := rfr.referralTrxRepo.GetMilestone(c, int64(cif))

	if err != nil {
		requestLogger.Debug(models.ErrMilestone)

		return nil, models.ErrMilestone
	}

	stage := os.Getenv(`STAGE`)
	limitRewardCounter := os.Getenv(`LIMIT_REWARD_COUNTER`)
	milestone.Stages, _ = strconv.ParseInt(stage, 10, 64)
	milestone.LimitRewardCounter, _ = strconv.ParseInt(limitRewardCounter, 10, 64)
	milestone.LimitReward = milestone.Stages * milestone.LimitRewardCounter

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
