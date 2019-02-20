package usecase

import (
	"context"
	"time"

	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
)

type campaignUseCase struct {
	campaignRepo   campaigns.Repository
	contextTimeout time.Duration
}

// NewCampaignUseCase will create new an campaignUseCase object representation of campaigns.UseCase interface
func NewCampaignUseCase(a campaigns.Repository, timeout time.Duration) campaigns.UseCase {
	return &campaignUseCase{
		campaignRepo:   a,
		contextTimeout: timeout,
	}
}

func (a *campaignUseCase) CreateCampaign(c context.Context, m *models.Campaign) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.campaignRepo.CreateCampaign(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (a *campaignUseCase) UpdateCampaign(c context.Context, id int64, updateCampaign *models.UpdateCampaign) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.campaignRepo.UpdateCampaign(ctx, id, updateCampaign)
	if err != nil {
		return err
	}

	return nil
}

func (a *campaignUseCase) GetCampaign(c context.Context, name string, status string, startDate string, endDate string) ([]*models.Campaign, error) {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	listCampaign, err := a.campaignRepo.GetCampaign(ctx, name, status, startDate, endDate)
	if err != nil {
		return nil, err
	}

	return listCampaign, nil
}
