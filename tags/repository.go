package tags

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the tags repository contract
type Repository interface {
	CreateTag(echo.Context, *models.Tag, int64) error
	Delete(echo.Context, int64) error
}
