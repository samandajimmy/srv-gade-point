package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the referrals usecases
type RefUseCase interface {
	UPostCoreTrx(echo.Context, []models.CoreTrxPayload) ([]models.CoreTrxResponse, error)
	UCreateReferralCodes(echo.Context, models.RequestCreateReferral) (models.RespReferral, error)
	UGetReferralCodes(echo.Context, models.RequestReferralCodeUser) (models.ReferralCodeUser, error)
	UReferralCIFValidate(echo.Context, string) (models.ReferralCodes, error)
	UValidateReferrer(c echo.Context, pl models.PayloadValidator, campaign *models.Campaign) (models.SumIncentive, error)
	UGetPrefixActiveCampaignReferral(echo.Context) (models.PrefixResponse, error)
	UGetHistoryIncentive(c echo.Context, pl models.RequestHistoryIncentive) ([]models.ResponseHistoryIncentive, error)
}
