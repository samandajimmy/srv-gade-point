package rewards

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the rewards usecases
type UseCase interface {
	CreateReward(echo.Context, *models.Reward, int64) error
	DeleteByCampaign(echo.Context, int64) error
	Inquiry(echo.Context, *models.PayloadValidator) (models.RewardsInquiry, error)
	Payment(echo.Context, *models.RewardPayment) error
}
