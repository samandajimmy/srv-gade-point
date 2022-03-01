package referraltrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referral transactions repository contract
type Repository interface {
	RPostReferralTrx(echo.Context, models.ReferralTrx) error
	RGetMilestone(echo.Context, models.MilestonePayload) (*models.Milestone, error)
	RGetRanking(echo.Context, models.RankingPayload) ([]*models.Ranking, error)
	RGetRankingByReferralCode(echo.Context, string) (*models.Ranking, error)
}
