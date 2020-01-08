package referraltrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the referraltrxs usecases
type UseCase interface {
	GetMilestone(echo.Context, models.MilestonePayload) (*models.Milestone, error)
	GetRanking(echo.Context, models.RankingPayload) ([]models.Ranking, error)
}
