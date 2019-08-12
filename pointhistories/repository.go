package pointhistories

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the pointhistories repository contract
type Repository interface {
	GetUsers(echo.Context, map[string]interface{}) ([]models.PointHistory, error)
	CountUsers(echo.Context, map[string]interface{}) (string, error)
	CountUserPointHistory(echo.Context, map[string]interface{}) (string, error)
	GetUserPointHistory(echo.Context, map[string]interface{}) ([]models.PointHistory, error)
	GetUserPoint(echo.Context, string) (float64, error)
}
