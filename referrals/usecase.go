package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the referrals usecases
type UseCase interface {
	CreateReferralCodes(echo.Context, models.RequestCreateReferral) (models.ReferralCodes, error)
}
