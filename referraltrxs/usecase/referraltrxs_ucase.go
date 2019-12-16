package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referraltrxs"
	"github.com/labstack/echo"
	"os"
	"strconv"
)

type referralTrxUseCase struct {
	referralTrxRepo referraltrxs.Repository
}

// NewReferralTrxUseCase will create new an referralTrxUseCase object representation of referraltrxs.UseCase interface
func NewReferralTrxUseCase(referralTrxRepo referraltrxs.Repository) referraltrxs.UseCase {
	return &referralTrxUseCase{
		referralTrxRepo: referralTrxRepo,
	}
}

func (rfr *referralTrxUseCase) GetMilestone(c echo.Context, CIF string) (*models.Milestone, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	cif, err := strconv.Atoi(CIF)

	if err != nil {
		requestLogger.Debug(err)

		return nil, models.ErrCIF
	}

	milestone, err := rfr.referralTrxRepo.GetMilestone(c, int64(cif))

	if err != nil {
		requestLogger.Debug(models.ErrMilestone)

		return nil, models.ErrMilestone
	}

	price := os.Getenv(`PRICE`)
	totalLimit := os.Getenv(`TOTAL_LIMIT`)
	milestone.Price, _ = strconv.ParseInt(price, 10, 64)
	milestone.Total, _ = strconv.ParseInt(totalLimit, 10, 64)
	milestone.TotalPrice = milestone.Price * milestone.Total

	return milestone, nil
}
