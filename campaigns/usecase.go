package campaigns

import (
	"context"

	"gade/srv-gade-point/models"
)

// UseCase represent the campaign's usecases
type UseCase interface {
	CreateCampaign(context.Context, *models.Campaign) error
	UpdateCampaign(ctx context.Context, id int64, updateCampaign *models.UpdateCampaign) (*models.Response, error)
	GetCampaign(ctx context.Context, name string, status string, startDate string, endDate string) ([]*models.Campaign, error)
}
