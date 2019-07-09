package quotas

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the quotas usecases
type UseCase interface {
	Create(echo.Context, *models.Quota, int64) error
	DeleteByReward(echo.Context, int64) error
}
