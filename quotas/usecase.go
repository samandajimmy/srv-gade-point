package quotas

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the quotas usecases
type UseCase interface {
	Create(echo.Context, *models.Quota, *models.Reward) error
	DeleteByReward(echo.Context, int64) error
	CheckQuota(echo.Context, models.Reward, string) (bool, error)
	UpdateAddQuota(echo.Context, int64) error
	UpdateReduceQuota(echo.Context, int64) error
}
