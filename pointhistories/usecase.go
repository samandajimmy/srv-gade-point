package pointhistories

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the pointhistories usecases
type UseCase interface {
	GetUsers(echo.Context, map[string]interface{}) ([]models.PointHistory, string, error)
	GetUserPointHistory(echo.Context, map[string]interface{}) ([]models.PointHistory, string, error)
	GetUserPoint(echo.Context, string) (*models.UserPoint, error)
}
