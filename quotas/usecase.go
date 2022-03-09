package quotas

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the quotas usecases
type QUseCase interface {
	Create(echo.Context, *models.Quota, *models.Reward) error
	DeleteByReward(echo.Context, int64) error
	CheckQuota(echo.Context, models.Reward, *models.PayloadValidator) (bool, error)
	UpdateAddQuota(echo.Context, int64) error
	UpdateReduceQuota(echo.Context, int64) error
}
