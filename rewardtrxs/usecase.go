package rewardtrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the reward transactions usecases
type UseCase interface {
	GetByRefID(echo.Context, string) (models.RewardsInquiry, error)
	CountByCIF(echo.Context, models.Quota, models.Reward, string) (int64, error)
	GetRewardTrxs(echo.Context, *models.RewardsPayload) ([]models.RewardTrx, string, error)
}
