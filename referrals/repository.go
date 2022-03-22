package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referrals repository contract
type RefRepository interface {
	CreateReferral(echo.Context, models.ReferralCodes) (models.ReferralCodes, error)
	GetCampaignId(echo.Context, string) (int64, error)
	GetReferralByCif(echo.Context, string) (models.ReferralCodes, error)
	RSumRefIncentive(c echo.Context, promoCode string, reward models.Reward) (models.SumIncentive, error)
	RGetReferralCampaignMetadata(c echo.Context, pv models.PayloadValidator) (models.PrefixResponse, error)
}
