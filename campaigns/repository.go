package campaigns

import (
	"context"
	"gade/srv-gade-point/models"
)

// Repository represent the campaign's repository contract
type Repository interface {
	CreateCampaign(ctx context.Context, a *models.Campaign) error
	UpdateCampaign(ctx context.Context, id int64, updateCampaign *models.UpdateCampaign) error
	GetCampaign(ctx context.Context, name string, status string, startDate string, endDate string, page int, limit int) ([]*models.Campaign, error)
	GetValidatorCampaign(ctx context.Context, a *models.GetCampaignValue) (*models.Campaign, error)
	SavePoint(ctx context.Context, a *models.CampaignTrx) error
	GetUserPoint(ctx context.Context, UserId string) (float64, error)
	GetUserPointHistory(ctx context.Context, UserID string) ([]models.CampaignTrx, error)
	CountCampaign(ctx context.Context, name string, status string, startDate string, endDate string) (int, error)
}
