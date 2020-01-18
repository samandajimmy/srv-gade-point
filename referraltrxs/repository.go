package referraltrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referral transactions repository contract
type Repository interface {
	Create(echo.Context, models.ReferralTrx) error
	GetMilestone(echo.Context, models.MilestonePayload) (*models.Milestone, error)
	IsReferralTrxExist(c echo.Context, refTrx models.ReferralTrx) (int64, error)
	GetTotalGoldbackReferrer(c echo.Context, refTrx models.ReferralTrx) (float64, error)
	GetRanking(echo.Context, string) ([]models.Ranking, error)
	GetRankingByReferralCode(echo.Context, string) (*models.Ranking, error)
}
