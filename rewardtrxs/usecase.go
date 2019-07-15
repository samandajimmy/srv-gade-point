package rewardtrxs

import (
	"github.com/labstack/echo"
)

// UseCase represent the reward transactions usecases
type UseCase interface {
	Create(echo.Context, map[string]interface{}, int64) error
	UpdateSuccess(echo.Context, map[string]interface{}) error
	UpdateReject(echo.Context, map[string]interface{}) error
}
