package rewardtrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the reward transactions usecases
type UseCase interface {
	Create(echo.Context, models.PayloadValidator, int64, []models.RewardResponse) (models.RewardTrx, error)
	GetByRefID(echo.Context, string) (models.RewardsInquiry, error)
	CountByCIF(echo.Context, models.Quota, models.Reward, string) (int64, error)
	GetRewardTrxs(echo.Context, *models.RewardsPayload) ([]models.RewardTrx, string, error)
}
