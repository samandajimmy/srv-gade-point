package campaigns

import (
	"context"

	"gade/srv-gade-point/models"
)

// Usecase represent the campaign's usecases
type Usecase interface {
	CreateCampaign(context.Context, *models.Campaign) error
}
