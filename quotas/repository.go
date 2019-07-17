package quotas

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the quotas repository contract
type Repository interface {
	Create(echo.Context, *models.Quota, int64) error
	DeleteByReward(echo.Context, int64) error
	CheckQuota(echo.Context, int64) ([]*models.Quota, error)
	UpdateAddQuota(echo.Context, int64) error
	UpdateMinQuota(echo.Context, int64) error
}
