package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the referrals usecases
type RefUseCase interface {
	UPostCoreTrx(echo.Context, []models.CoreTrxPayload) ([]models.CoreTrxResponse, error)
	CreateReferralCodes(echo.Context, models.RequestCreateReferral) (models.ReferralCodes, error)
	GetReferralCodes(echo.Context, models.RequestReferralCodeUser) (models.ReferralCodeUser, error)
	UMaxIncentiveValidation(c echo.Context, refCode string) (models.SumIncentive, error)
	UReferralCIFValidate(echo.Context, string) (models.ReferralCodes, error)
	UValidateReferrer(c echo.Context, pl models.PayloadValidator, campaign *models.Campaign) (models.SumIncentive, error)
}
