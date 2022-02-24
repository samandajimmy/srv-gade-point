package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referrals repository contract
type Repository interface {
	PostCoreTrx(echo.Context, []models.CoreTrxPayload) error
}
