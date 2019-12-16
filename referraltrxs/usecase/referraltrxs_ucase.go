package usecase

import (
	"errors"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referraltrxs"
	"github.com/labstack/echo"
	"os"
	"strconv"
)

type referralTrxUseCase struct {
	referralTrxRepo referraltrxs.Repository
}

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

		return nil, errors.New("Something went wrong with input CIF")
	}

	milestone, err := rfr.referralTrxRepo.GetMilestone(c, int64(cif))

	price := os.Getenv(`PRICE`)
	totalLimit := os.Getenv(`TOTAL_LIMIT`)

	milestone.Price, err = strconv.ParseInt(price, 10, 64)
	milestone.Total, err = strconv.ParseInt(totalLimit, 10, 64)
	milestone.TotalPrice = milestone.Price * milestone.Total

	if err != nil {
		requestLogger.Debug(models.ErrMilestone)

		return nil, models.ErrMilestone
	}

	return milestone, nil
}
