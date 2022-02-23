package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referrals repository contract
type Repository interface {
	CreateReferralCodes(echo.Context, *models.ReferralCodes) error
	GetCampaignByPrefix(echo.Context, string) (int64, error)
	GetReferralCodesByCif(echo.Context, *models.ReferralCodes) error
}
