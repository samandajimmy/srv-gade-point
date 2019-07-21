package rewardtrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the reward transactions repository contract
type Repository interface {
	Create(echo.Context, models.PayloadValidator, int64, []models.RewardResponse) (models.RewardTrx, error)
	GetByRefID(echo.Context, string) (models.RewardsInquiry, error)
	CountByCIF(echo.Context, models.Quota, models.Reward, string) (int64, error)
	UpdateRewardTrx(echo.Context, *models.RewardPayment, int64) error
	CheckTrx(echo.Context, string, string, string) (*models.RewardTrx, error)
}
