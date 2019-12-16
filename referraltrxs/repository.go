package referraltrxs

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referral transactions repository contract
type Repository interface {
	Create(echo.Context, models.ReferralTrx) error
	CheckReferralTrxByExistingReferrer(c echo.Context, refTrx models.ReferralTrx) (int64, error)
	CheckReferralTrxByValueRewards(c echo.Context, refTrx models.ReferralTrx) (*models.ReferralTrx, error)
}
