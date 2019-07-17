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
) quotas.UseCase {
	return &quotaUseCase{
		quotaRepo: quotRepo,
		rwdTrxUC:  rwdTrxUC,
	}
}

func (quot *quotaUseCase) Create(c echo.Context, quota *models.Quota, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := quot.quotaRepo.Create(c, quota, rewardID)

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

func (quot *quotaUseCase) CheckQuota(c echo.Context, rewardID int64, cif string) (bool, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	quotas, err := quot.quotaRepo.CheckQuota(c, rewardID)

	if err != nil {
		requestLogger.Debug(models.ErrCheckQuotaFailed)

		return false, models.ErrCheckQuotaFailed
	}

	// validate all quota
	for _, quota := range quotas {

		// check quota stock
		if *quota.Available <= models.IsLimitAmount {
			requestLogger.Debug(models.ErrQuotaNotAvailable)

			return false, models.ErrQuotaNotAvailable
		}

		// check quota allocation
		if *quota.IsPerUser == models.IsPerUserTrue {

			count, err := quot.rwdTrxUC.CountByCIF(c, *quota, cif)

			if err != nil {
				requestLogger.Debug(models.ErrCheckQuotaFailed)

				return false, models.ErrCheckQuotaFailed
			}

			if count < *quota.Available {
				requestLogger.Debug(models.ErrQuotaNotAvailable)

				return false, models.ErrQuotaNotAvailable
			}
		}

	}

	return true, nil
}
