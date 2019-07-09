package quotas

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the quotas repository contract
type Repository interface {
	Create(echo.Context, *models.Quota, int64) error
	DeleteByReward(echo.Context, int64) error
}
