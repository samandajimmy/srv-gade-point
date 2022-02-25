package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referrals repository contract
type Repository interface {
	PostCoreTrx(echo.Context, []models.CoreTrxPayload) error
	CreateReferral(echo.Context, models.ReferralCodes) (models.ReferralCodes, error)
	GetCampaignId(echo.Context, string) (int64, error)
	GetReferralByCif(echo.Context, string) (models.ReferralCodes, error)
}
