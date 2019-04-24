package campaigntrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the campaigntrx's usecases
type UseCase interface {
	GetUsers(echo.Context, map[string]interface{}) ([]models.CampaignTrx, string, error)
}
