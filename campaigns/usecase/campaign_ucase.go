package usecase

import (
	"context"
	"time"

	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
)

type campaignUsecase struct {
	campaignRepo   campaigns.Repository
	contextTimeout time.Duration
}

// NewCampaignUsecase will create new an campaignUsecase object representation of campaigns.Usecase interface
func NewCampaignUsecase(a campaigns.Repository, timeout time.Duration) campaigns.Usecase {
	return &campaignUsecase{
		campaignRepo:   a,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */

func (a *campaignUsecase) CreateCampaign(c context.Context, m *models.Campaign) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.campaignRepo.CreateCampaign(ctx, m)
	if err != nil {
		return err
	}
	return nil
}
