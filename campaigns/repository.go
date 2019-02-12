package campaigns

import (
	"context"

	"gade/srv-gade-point/models"
)

// Repository represent the campaign's repository contract
type Repository interface {
	CreateCampaign(context.Context, *models.Campaign) error
}
