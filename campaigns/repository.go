package campaigns

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the campaign's repository contract
type CRepository interface {
	CreateCampaign(echo.Context, *models.Campaign) error
	UpdateCampaign(echo.Context, int64, *models.Campaign) error
	GetCampaignDetail(echo.Context, int64) (*models.Campaign, error)
	GetCampaign(echo.Context, map[string]interface{}) ([]*models.Campaign, error)
	SavePoint(echo.Context, *models.CampaignTrx) error
	CountCampaign(echo.Context, map[string]interface{}) (int, error)
	UpdateExpiryDate(echo.Context) error
	UpdateStatusBasedOnStartDate() error
	GetCampaignAvailable(echo.Context, models.PayloadValidator) ([]*models.Campaign, error)
	GetReferralCampaign(c echo.Context, pv models.PayloadValidator) *[]*models.Campaign
	Delete(echo.Context, int64) error
}
