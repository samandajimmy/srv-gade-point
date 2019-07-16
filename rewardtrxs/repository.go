package rewardtrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the reward transactions repository contract
type Repository interface {
	Create(echo.Context, models.PayloadValidator, int64, []models.RewardResponse) (models.RewardTrx, error)
	UpdateSuccess(echo.Context, map[string]interface{}) error
	UpdateReject(echo.Context, map[string]interface{}) error
	GetByRefID(echo.Context, string) (models.RewardsInquiry, error)
}
