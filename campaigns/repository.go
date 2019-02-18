package campaigns

import (
	"context"

	"gade/srv-gade-point/models"
)

// Repository represent the campaign's repository contract
type Repository interface {
	CreateCampaign(ctx context.Context, a *models.Campaign) error
	UpdateCampaign(ctx context.Context, id int64, updateCampaign *models.UpdateCampaign) error
	GetCampaign(ctx context.Context, name string, status string, startDate string, endDate string) ([]*models.Campaign, error)
}
