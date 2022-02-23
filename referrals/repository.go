package referrals

import (
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

// Repository represent the referrals repository contract
type Repository interface {
	CreateCoreTrx(echo.Context, models.CoreTrxPayload) error
}
