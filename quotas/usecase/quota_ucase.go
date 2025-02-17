package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/quotas"
	"gade/srv-gade-point/rewardtrxs"

	"github.com/labstack/echo"
)

type quotaUseCase struct {
	quotaRepo quotas.Repository
	rwdTrxUC  rewardtrxs.UseCase
}

// NewQuotaUseCase will create new an quotaUseCase object representation of quotas.UseCase interface
func NewQuotaUseCase(
	quotRepo quotas.Repository,
	rwdTrxUC rewardtrxs.UseCase,
) quotas.QUseCase {
	return &quotaUseCase{
		quotaRepo: quotRepo,
		rwdTrxUC:  rwdTrxUC,
	}
}

func (quot *quotaUseCase) Create(c echo.Context, quota *models.Quota, reward *models.Reward) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := quot.quotaRepo.Create(c, quota, reward)

	if err != nil {
		requestLogger.Debug(models.ErrQuotaFailed)

		return models.ErrQuotaFailed
	}

	return nil
}

func (quot *quotaUseCase) DeleteByReward(c echo.Context, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := quot.quotaRepo.DeleteByReward(c, rewardID)

	if err != nil {
		requestLogger.Debug(models.ErrDelQuotaFailed)

		return models.ErrDelQuotaFailed
	}

	return nil
}

func (quot *quotaUseCase) refreshQuota(c echo.Context, reward models.Reward, plValidator *models.PayloadValidator) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// Check Refresh Quota List
	quotaList, err := quot.quotaRepo.CheckRefreshQuota(c, plValidator)

	if err != nil {
		requestLogger.Debug(models.ErrCheckQuotaFailed)

		return models.ErrCheckQuotaFailed
	}

	for _, qList := range quotaList {
		err := quot.quotaRepo.RefreshQuota(c, qList, plValidator)

		if err != nil {
			requestLogger.Debug(models.ErrRefreshQuotaFailed)

			return models.ErrRefreshQuotaFailed
		}
	}

	return nil
}

func (quot *quotaUseCase) CheckQuota(c echo.Context, reward models.Reward, plValidator *models.PayloadValidator) (bool, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// Check Refresh Quota List
	if err := quot.refreshQuota(c, reward, plValidator); err != nil {
		requestLogger.Debug(models.ErrRefreshQuotaFailed)

		return false, models.ErrRefreshQuotaFailed
	}

	// Check quota
	quotas, err := quot.quotaRepo.CheckQuota(c, reward.ID)

	if err != nil {
		requestLogger.Debug(models.ErrCheckQuotaFailed)

		return false, models.ErrCheckQuotaFailed
	}

	// validate all quota
	for _, quota := range quotas {
		// check quota allocation
		if *quota.IsPerUser == models.IsPerUserFalse && *quota.Amount == models.QuotaUnlimited {
			return true, nil
		}

		if *quota.IsPerUser == models.IsPerUserTrue && plValidator.CIF != "" {
			count, err := quot.rwdTrxUC.CountByCIF(c, *quota, reward, plValidator.CIF)

			if err != nil {
				requestLogger.Debug(models.ErrCheckQuotaFailed)

				return false, models.ErrCheckQuotaFailed
			}

			if count >= *quota.Available {
				requestLogger.Debug(models.ErrQuotaNACIF)

				return false, models.ErrQuotaNACIF
			}
		}

		// check quota availibility
		if quota.Available != nil && *quota.Available <= models.IsLimitAmount {
			err = models.ErrTodaysQuotaNotAvailable

			if *quota.NumberOfDays == 0 {
				err = models.ErrQuotaNotAvailable
			}

			requestLogger.Debug(err)

			return false, err
		}
	}

	return true, nil
}

func (quot *quotaUseCase) UpdateAddQuota(c echo.Context, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	if rewardID == 0 {
		return nil
	}

	err := quot.quotaRepo.UpdateAddQuota(c, rewardID)

	if err != nil {
		requestLogger.Debug(models.ErrAddQuotaFailed)

		return models.ErrAddQuotaFailed
	}

	return nil
}

func (quot *quotaUseCase) UpdateReduceQuota(c echo.Context, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	if rewardID == 0 {
		return nil
	}

	err := quot.quotaRepo.UpdateReduceQuota(c, rewardID)

	if err != nil {
		requestLogger.Debug(models.ErrReduceQuotaFailed)

		return models.ErrReduceQuotaFailed
	}

	return nil
}
