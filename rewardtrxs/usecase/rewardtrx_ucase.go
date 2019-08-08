package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewardtrxs"
	"strconv"

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

func (rwdTrx *rewardTrxUseCase) GetByRefID(c echo.Context, refID string) (models.RewardsInquiry, error) {
	var rewardTrx models.RewardsInquiry
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	rewardTrx, err := rwdTrx.rewardTrxRepo.GetByRefID(c, refID)

	if err != nil {
		requestLogger.Debug(models.ErrRefTrxNotFound)

		return rewardTrx, models.ErrRefTrxNotFound
	}

	return rewardTrx, nil
}

func (rwdTrx *rewardTrxUseCase) CountByCIF(c echo.Context, quot models.Quota, reward models.Reward, cif string) (int64, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	count, err := rwdTrx.rewardTrxRepo.CountByCIF(c, quot, reward, cif)

	if err != nil {
		requestLogger.Debug(models.ErrCheckQuotaFailed)

		return 0, models.ErrCheckQuotaFailed
	}

	return count, nil
}

func (rwdTrx *rewardTrxUseCase) GetRewardTrxs(c echo.Context, rewardPayload *models.RewardsPayload) ([]models.RewardTrx, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// counter data
	counter, err := rwdTrx.rewardTrxRepo.CountRewardTrxs(c, rewardPayload)
	if err != nil {
		requestLogger.Debug(models.ErrGetRewardTrxCounter)

		return nil, "", models.ErrGetRewardTrxCounter
	}

	// get data
	data, err := rwdTrx.rewardTrxRepo.GetRewardTrxs(c, rewardPayload)

	if err != nil {
		requestLogger.Debug(models.ErrGetRewardTrx)

		return nil, "", models.ErrGetRewardTrx
	}

	return data, strconv.FormatInt(counter, 10), err
}
