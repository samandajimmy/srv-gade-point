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

func (rwdTrx *rewardTrxUseCase) Create(c echo.Context, payload map[string]interface{}, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := rwdTrx.rewardTrxRepo.Create(c, payload, rewardID)

	if err != nil {
		requestLogger.Debug(models.ErrRewardTrxFailed)

		return models.ErrRewardTrxFailed
	}

	return nil
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
