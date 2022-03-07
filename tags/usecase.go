package tags

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the tags usecases
type TUseCase interface {
	CreateTag(echo.Context, *models.Tag, int64) error
	Delete(echo.Context, int64) error
}
