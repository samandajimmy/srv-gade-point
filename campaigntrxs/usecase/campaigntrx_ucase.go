package usecase

import (
	"gade/srv-gade-point/campaigntrxs"
	"gade/srv-gade-point/models"

	"github.com/labstack/echo"
)

type campaignTrxUseCase struct {
	campaignTrxRepo campaigntrxs.Repository
}

// NewCampaignTrxUseCase will create new an campaignTrxUseCase object representation of campaigntrxs.UseCase interface
func NewCampaignTrxUseCase(cmpgnTrxRepo campaigntrxs.Repository) campaigntrxs.UseCase {
	return &campaignTrxUseCase{
		campaignTrxRepo: cmpgnTrxRepo,
	}
}

func (cmpgnTrxUs *campaignTrxUseCase) GetUsers(c echo.Context, payload map[string]interface{}) ([]models.CampaignTrx, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	// check users availabilty
	counter, err := cmpgnTrxUs.campaignTrxRepo.CountUsers(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrUsersNA)

		return nil, "", models.ErrUsersNA
	}

	// get users point data
	data, err := cmpgnTrxUs.campaignTrxRepo.GetUsers(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetUsersPoint)

		return nil, "", models.ErrGetUsersPoint
	}

	return data, counter, nil
}
