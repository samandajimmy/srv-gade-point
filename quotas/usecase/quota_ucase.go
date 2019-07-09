package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/quotas"

	"github.com/labstack/echo"
)

type quotaUseCase struct {
	quotaRepo quotas.Repository
}

// NewQuotaUseCase will create new an quotaUseCase object representation of quotas.UseCase interface
func NewQuotaUseCase(quotRepo quotas.Repository) quotas.UseCase {
	return &quotaUseCase{
		quotaRepo: quotRepo,
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
