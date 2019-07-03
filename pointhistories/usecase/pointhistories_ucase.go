package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/pointhistories"

	"github.com/labstack/echo"
)

type pointHistoryUseCase struct {
	pointHistoryRepo pointhistories.Repository
}

// NewPointHistoryUseCase will create new an pointHistoryUseCase object representation of pointhistories.UseCase interface
func NewPointHistoryUseCase(pntHstryRepo pointhistories.Repository) pointhistories.UseCase {
	return &pointHistoryUseCase{
		pointHistoryRepo: pntHstryRepo,
	}
}

func (pntHstryUs *pointHistoryUseCase) GetUsers(c echo.Context, payload map[string]interface{}) ([]models.PointHistory, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// check users availabilty
	counter, err := pntHstryUs.pointHistoryRepo.CountUsers(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrUsersNA)

		return nil, "", models.ErrUsersNA
	}

	// get users point data
	data, err := pntHstryUs.pointHistoryRepo.GetUsers(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetUsersPoint)

		return nil, "", models.ErrGetUsersPoint
	}

	return data, counter, nil
}

func (pntHstryUs *pointHistoryUseCase) GetUserPointHistory(c echo.Context, payload map[string]interface{}) ([]models.PointHistory, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	counter, err := pntHstryUs.pointHistoryRepo.CountUserPointHistory(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrUserPointHistoryNA)

		return nil, "", models.ErrUserPointHistoryNA
	}

	dataHistory, err := pntHstryUs.pointHistoryRepo.GetUserPointHistory(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetUserPointHistory)

		return nil, "", models.ErrGetUserPointHistory
	}

	return dataHistory, counter, nil
}

func (pntHstryUs *pointHistoryUseCase) GetUserPoint(c echo.Context, CIF string) (*models.UserPoint, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	p := new(models.UserPoint)
	zero := float64(0)
	pointAmount, err := pntHstryUs.pointHistoryRepo.GetUserPoint(c, CIF)

	if err != nil {
		requestLogger.Debug(models.ErrGetUserPoint)
		p.UserPoint = &zero

		return p, models.ErrGetUserPoint
	}

	if pointAmount == 0 {
		requestLogger.Debug(models.ErrUserPointNA)
		p.UserPoint = &zero

		return p, models.ErrUserPointNA
	}

	p.UserPoint = &pointAmount

	return p, nil
}
