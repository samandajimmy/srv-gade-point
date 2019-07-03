package campaigns

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the campaign's usecases
type UseCase interface {
	CreateCampaign(echo.Context, *models.Campaign) error
	UpdateCampaign(echo.Context, string, *models.UpdateCampaign) error
	GetCampaignDetail(echo.Context, string) (*models.Campaign, error)
	GetCampaign(echo.Context, map[string]interface{}) (string, []*models.Campaign, error)
	GetCampaignValue(echo.Context, *models.GetCampaignValue) (*models.UserPoint, error)
	UpdateStatusBasedOnStartDate() error
}
