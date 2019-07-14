package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/tags"

	"github.com/labstack/echo"
)

type tagUseCase struct {
	tagRepo tags.Repository
}

// NewTagUseCase will create new an tagUseCase object representation of tags.UseCase interface
func NewTagUseCase(tgRepo tags.Repository) tags.UseCase {
	return &tagUseCase{
		tagRepo: tgRepo,
	}
}

func (tg *tagUseCase) CreateTag(c echo.Context, tag *models.Tag, campaignID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := tg.tagRepo.CreateTag(c, tag, campaignID)

	if err != nil {
		requestLogger.Debug(models.ErrTagFailed)

		return models.ErrTagFailed
	}

	return nil
}

func (tg *tagUseCase) Delete(c echo.Context, tagID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := tg.tagRepo.Delete(c, tagID)

	if err != nil {
		requestLogger.Debug(models.ErrDelQuotaFailed)

		return models.ErrDelQuotaFailed
	}

	return nil
}
