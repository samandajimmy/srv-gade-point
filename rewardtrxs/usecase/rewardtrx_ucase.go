package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewardtrxs"

	"github.com/labstack/echo"
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

func (rwdTrx *rewardTrxUseCase) Create(c echo.Context, payload models.PayloadValidator, rewardID int64, resp []models.RewardResponse) (models.RewardTrx, error) {
	var rewardTrx models.RewardTrx
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	rewardTrx, err := rwdTrx.rewardTrxRepo.Create(c, payload, rewardID, resp)

	if err != nil {
		requestLogger.Debug(models.ErrRewardTrxFailed)

		return rewardTrx, models.ErrRewardTrxFailed
	}

	return rewardTrx, nil
}

func (rwdTrx *rewardTrxUseCase) UpdateSuccess(c echo.Context, payload map[string]interface{}) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := rwdTrx.rewardTrxRepo.UpdateSuccess(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrRewardTrxUpdateFailed)

		return models.ErrRewardTrxUpdateFailed
	}

	return nil
}

func (rwdTrx *rewardTrxUseCase) UpdateReject(c echo.Context, payload map[string]interface{}) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := rwdTrx.rewardTrxRepo.UpdateReject(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrRewardTrxUpdateFailed)

		return models.ErrRewardTrxUpdateFailed
	}

	return nil
}

func (rwdTrx *rewardTrxUseCase) GetByRefID(c echo.Context, refID string) (models.RewardTrx, error) {
	var rewardTrx models.RewardTrx
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	rewardTrx, err := rwdTrx.rewardTrxRepo.GetByRefID(c, refID)

	if err != nil {
		requestLogger.Debug(models.ErrRewardTrxUpdateFailed)

		return rewardTrx, models.ErrRewardTrxUpdateFailed
	}

	return rewardTrx, nil
}
