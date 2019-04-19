package campaigns

import (
	"context"
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the campaign's repository contract
type Repository interface {
	CreateCampaign(echo.Context, *models.Campaign) error
	UpdateCampaign(echo.Context, int64, *models.UpdateCampaign) error
	GetCampaignDetail(echo.Context, int64) (*models.Campaign, error)
	GetCampaign(ctx context.Context, name string, status string, startDate string, endDate string, page int, limit int) ([]*models.Campaign, error)
	GetValidatorCampaign(ctx context.Context, a *models.GetCampaignValue) (*models.Campaign, error)
	SavePoint(ctx context.Context, a *models.CampaignTrx) error
	GetUserPoint(ctx context.Context, UserID string) (float64, error)
	GetUserPointHistory(ctx context.Context, UserID string) ([]models.CampaignTrx, error)
	CountCampaign(ctx context.Context, name string, status string, startDate string, endDate string) (int, error)
	UpdateExpiryDate(ctx context.Context) error
	UpdateStatusBasedOnStartDate() error
}
