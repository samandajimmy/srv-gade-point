package referraltrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

type UseCase interface {
	GetMilestone(echo.Context, string) (*models.Milestone, error)
}
