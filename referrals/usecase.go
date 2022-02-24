package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// UseCase represent the referrals usecases
type UseCase interface {
	PostCoreTrx(echo.Context, []models.CoreTrxPayload) ([]models.CoreTrxResponse, error)
}
