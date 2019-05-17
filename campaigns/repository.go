package campaigns

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the campaign's repository contract
type Repository interface {
	CreateCampaign(echo.Context, *models.Campaign) error
	UpdateCampaign(echo.Context, int64, *models.UpdateCampaign) error
	GetCampaignDetail(echo.Context, int64) (*models.Campaign, error)
	GetCampaign(echo.Context, map[string]interface{}) ([]*models.Campaign, error)
	GetValidatorCampaign(echo.Context, *models.GetCampaignValue) (*models.Campaign, error)
	SavePoint(echo.Context, *models.CampaignTrx) error
	GetUserPoint(echo.Context, string) (float64, error)
	GetUserPointHistory(echo.Context, map[string]interface{}) ([]models.CampaignTrx, error)
	CountUserPointHistory(echo.Context, map[string]interface{}) (string, error)
	CountCampaign(echo.Context, map[string]interface{}) (int, error)
	UpdateExpiryDate(echo.Context) error
	UpdateStatusBasedOnStartDate() error
	GetCampaignAvailable(echo.Context) ([]*models.Campaign, error)
}
