package rewardtrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the reward transactions usecases
type UseCase interface {
	Create(echo.Context, models.PayloadValidator, int64) error
	UpdateSuccess(echo.Context, map[string]interface{}) error
	UpdateReject(echo.Context, map[string]interface{}) error
}
