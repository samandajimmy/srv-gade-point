package campaigns

import (
	"context"
	"gade/srv-gade-point/models"
)

// UseCase represent the campaign's usecases
type UseCase interface {
	CreateCampaign(context.Context, *models.Campaign) error
	UpdateCampaign(ctx context.Context, id int64, updateCampaign *models.UpdateCampaign) error
	GetCampaign(ctx context.Context, name string, status string, startDate string, endDate string, page int, limit int) (string, []*models.Campaign, error)
	GetCampaignDetail(ctx context.Context, id int64) (*models.Campaign, error)
	GetCampaignValue(context.Context, *models.GetCampaignValue) (*models.UserPoint, error)
	GetUserPoint(ctx context.Context, userID string) (*models.UserPoint, error)
	GetUserPointHistory(ctx context.Context, userID string) ([]models.CampaignTrx, error)
	UpdateStatusBasedOnStartDate() error
}
