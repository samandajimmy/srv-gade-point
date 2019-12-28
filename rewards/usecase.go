package rewards

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the rewards usecases
type UseCase interface {
	CreateReward(echo.Context, *models.Reward, *models.Campaign) error
	DeleteByCampaign(echo.Context, int64) error
	Inquiry(echo.Context, *models.PayloadValidator) (models.RewardsInquiry, *models.ResponseErrors)
	Payment(echo.Context, *models.RewardPayment) (models.RewardTrxResponse, error)
	CheckTransaction(echo.Context, *models.RewardPayment) (models.RewardTrxResponse, error)
	RefreshTrx()
	GetRewards(echo.Context, *models.RewardsPayload) ([]models.Reward, string, error)
	GetRewardPromotions(echo.Context, models.PayloadValidator) ([]*models.RewardPromotions, *models.ResponseErrors, string, error)
}
