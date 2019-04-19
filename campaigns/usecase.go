package campaigns

import (
	"context"
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the campaign's usecases
type UseCase interface {
	CreateCampaign(echo.Context, *models.Campaign) error
	UpdateCampaign(echo.Context, string, *models.UpdateCampaign) error
	GetCampaignDetail(echo.Context, string) (*models.Campaign, error)
	GetCampaign(ctx context.Context, name string, status string, startDate string, endDate string, page int, limit int) (string, []*models.Campaign, error)
	GetCampaignValue(context.Context, *models.GetCampaignValue) (*models.UserPoint, error)
	GetUserPoint(ctx context.Context, userID string) (*models.UserPoint, error)
	GetUserPointHistory(ctx context.Context, userID string) ([]models.CampaignTrx, error)
	UpdateStatusBasedOnStartDate() error
}
