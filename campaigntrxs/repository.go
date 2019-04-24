package campaigntrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the campaigntrx's repository contract
type Repository interface {
	GetUsers(echo.Context, map[string]interface{}) ([]models.CampaignTrx, error)
	CountUsers(echo.Context, map[string]interface{}) (string, error)
}
