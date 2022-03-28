package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referrals repository contract
type RefRepository interface {
	RPostCoreTrx(echo.Context, []models.CoreTrxPayload) error
	RCreateReferral(echo.Context, models.ReferralCodes) (models.ReferralCodes, error)
	RGetCampaignId(echo.Context, string) (int64, error)
	RGetReferralByCif(echo.Context, string) (models.ReferralCodes, error)
	RSumRefIncentive(c echo.Context, promoCode string, reward models.Reward) (models.SumIncentive, error)
	RGenerateCode(c echo.Context, refCode models.ReferralCodes, prefix string) string
	RGetReferralCampaignMetadata(c echo.Context, pv models.PayloadValidator) (models.PrefixResponse, error)
	RGetHistoryIncentive(c echo.Context, refCif string) ([]models.ResponseHistoryIncentive, error)
	RFriendsReferral(echo.Context, models.PayloadFriends) ([]models.RespFriends, error)
}
