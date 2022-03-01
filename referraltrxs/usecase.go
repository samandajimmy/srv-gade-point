package referraltrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the referraltrxs usecases
type UseCase interface {
	UGetMilestone(echo.Context, models.MilestonePayload) (*models.Milestone, error)
	UGetRanking(echo.Context, models.RankingPayload) ([]*models.Ranking, error)
}
