package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the referrals usecases
type UseCase interface {
	UPostCoreTrx(echo.Context, []models.CoreTrxPayload) ([]models.CoreTrxResponse, error)
	CreateReferralCodes(echo.Context, models.RequestCreateReferral) (models.ReferralCodes, error)
	GetReferralCodes(echo.Context, models.RequestReferralCodeUser) (models.ReferralCodeUser, error)
	UMaxIncentiveValidation(c echo.Context, refCode string) (models.SumIncentive, error)
}
