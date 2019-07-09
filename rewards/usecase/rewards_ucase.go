package usecase

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/quotas"
	"gade/srv-gade-point/rewards"
	"gade/srv-gade-point/tags"

	"github.com/labstack/echo"
)

type rewardUseCase struct {
	rewardRepo rewards.Repository
	tagUC      tags.UseCase
	quotaUC    quotas.UseCase
}

// NewRewardUseCase will create new an rewardUseCase object representation of rewards.UseCase interface
func NewRewardUseCase(rwdRepo rewards.Repository, tagUC tags.UseCase, quotaUC quotas.UseCase) rewards.UseCase {
	return &rewardUseCase{
		rewardRepo: rwdRepo,
		tagUC:      tagUC,
		quotaUC:    quotaUC,
	}
}

func (rwd *rewardUseCase) CreateReward(c echo.Context, reward *models.Reward, campaignID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := rwd.rewardRepo.CreateReward(c, reward, campaignID)

	if err != nil {
		requestLogger.Debug(models.ErrRewardFailed)

		return models.ErrRewardFailed
	}

	// create array quotas
	for _, quota := range *reward.Quotas {
		err = rwd.quotaUC.Create(c, &quota, reward.ID)

		if err != nil {
			break
		}
	}

	if err != nil {
		_ = rwd.quotaUC.DeleteByReward(c, reward.ID)
		requestLogger.Debug(models.ErrCreateQuotasFailed)

		return models.ErrCreateQuotasFailed
	}

	// create array tags
	for _, tag := range *reward.Tags {
		err = rwd.tagUC.CreateTag(c, &tag, reward.ID)

		if err != nil {
			break
		}
	}

	if err != nil {
		_ = rwd.tagUC.DeleteByReward(c, reward.ID)
		requestLogger.Debug(models.ErrCreateTagsFailed)

		return models.ErrCreateTagsFailed
	}

	return nil
}

func (rwd *rewardUseCase) DeleteByCampaign(c echo.Context, campaignID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := rwd.rewardRepo.DeleteByCampaign(c, campaignID)

	if err != nil {
		requestLogger.Debug(models.ErrDelRewardFailed)

		return models.ErrDelRewardFailed
	}

	return nil
}
