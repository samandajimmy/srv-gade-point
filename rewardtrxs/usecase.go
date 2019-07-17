package rewardtrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the reward transactions usecases
type UseCase interface {
	Create(echo.Context, models.PayloadValidator, int64, []models.RewardResponse) (models.RewardTrx, error)
	UpdateSuccess(echo.Context, map[string]interface{}) error
	UpdateReject(echo.Context, map[string]interface{}) error
	GetByRefID(echo.Context, string) (models.RewardsInquiry, error)
	CountByCIF(echo.Context, models.Quota, string) (int64, error)
}
