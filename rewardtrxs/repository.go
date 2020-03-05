package rewardtrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the reward transactions repository contract
type Repository interface {
	Create(echo.Context, models.PayloadValidator, models.RewardsInquiry) ([]*models.RewardTrx, error)
	GetByRefID(echo.Context, string) (models.RewardsInquiry, error)
	CountByCIF(echo.Context, models.Quota, models.Reward, string) (int64, error)
	UpdateRewardTrx(echo.Context, *models.RewardPayment, int64) error
	CheckTrx(echo.Context, string) (*models.RewardTrx, error)
	CheckRefID(echo.Context, string) (*models.RewardTrx, error)
	GetRewardByPayload(echo.Context, models.PayloadValidator) ([]*models.Reward, error)
	RewardTrxTimeout(models.RewardTrx)
	UpdateTimeoutTrx() error
	GetInquiredTrx() ([]models.RewardTrx, error)
	GetRewardTrxs(echo.Context, *models.RewardsPayload) ([]models.RewardTrx, error)
	CountRewardTrxs(echo.Context, *models.RewardsPayload) (int64, error)
	CheckRootRefId(c echo.Context, rootRefID string) ([]string, error)
}
