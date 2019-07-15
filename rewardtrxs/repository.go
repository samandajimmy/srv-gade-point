package rewardtrxs

import (
	"github.com/labstack/echo"
)

// Repository represent the reward transactions repository contract
type Repository interface {
	Create(echo.Context, map[string]interface{}, int64) error
	UpdateSuccess(echo.Context, map[string]interface{}) error
	UpdateReject(echo.Context, map[string]interface{}) error
}
