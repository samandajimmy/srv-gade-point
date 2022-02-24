package referraltrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referral transactions repository contract
type Repository interface {
	PostReferralTrx(echo.Context, models.ReferralTrx) error
	GetMilestone(echo.Context, models.MilestonePayload) (*models.Milestone, error)
	GetRanking(echo.Context, models.RankingPayload) ([]*models.Ranking, error)
	GetRankingByReferralCode(echo.Context, string) (*models.Ranking, error)
}
