package usecase

import (
	"context"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/labstack/gommon/log"
	govaluate "gopkg.in/Knetic/govaluate.v2"
)

var floatType = reflect.TypeOf(float64(0))

type campaignUseCase struct {
	campaignRepo   campaigns.Repository
	contextTimeout time.Duration
}

// NewCampaignUseCase will create new an campaignUseCase object representation of campaigns.UseCase interface
func NewCampaignUseCase(cmpgn campaigns.Repository, timeout time.Duration) campaigns.UseCase {
	return &campaignUseCase{
		campaignRepo:   cmpgn,
		contextTimeout: timeout,
	}
}

func (cmpgn *campaignUseCase) CreateCampaign(c context.Context, m *models.Campaign) error {
	ctx, cancel := context.WithTimeout(c, cmpgn.contextTimeout)
	defer cancel()

	err := cmpgn.campaignRepo.CreateCampaign(ctx, m)
	if err != nil {
		return err
	}
	return nil
}

func (cmpgn *campaignUseCase) UpdateCampaign(c context.Context, id int64, updateCampaign *models.UpdateCampaign) error {
	ctx, cancel := context.WithTimeout(c, cmpgn.contextTimeout)
	defer cancel()

	err := cmpgn.campaignRepo.UpdateCampaign(ctx, id, updateCampaign)
	if err != nil {
		return err
	}

	return nil
}

func (cmpgn *campaignUseCase) GetCampaign(c context.Context, name string, status string, startDate string, endDate string, page int, limit int) (string, []*models.Campaign, error) {
	ctx, cancel := context.WithTimeout(c, cmpgn.contextTimeout)
	defer cancel()

	listCampaign, err := cmpgn.campaignRepo.GetCampaign(ctx, name, status, startDate, endDate, page, limit)

	if err != nil {
		return "", nil, err
	}

	countCampaign, err := cmpgn.campaignRepo.CountCampaign(ctx, name, status, startDate, endDate)

	if err != nil {
		return "", nil, err
	}

	if countCampaign <= 0 {
		return "", listCampaign, nil
	}

	return strconv.Itoa(countCampaign), listCampaign, nil
}

func (cmpgn *campaignUseCase) GetCampaignValue(c context.Context, m *models.GetCampaignValue) (*models.UserPoint, error) {
	payloadValidator := &models.PayloadValidator{}
	payloadValidator.Validators = &models.Validator{}
	now := time.Now()
	ctx, cancel := context.WithTimeout(c, cmpgn.contextTimeout)
	defer cancel()

	// get available campaign
	campaigns, err := cmpgn.campaignRepo.GetCampaignAvailable(ctx)

	if err != nil {
		// no campaign available
		log.Error(err)

		return nil, models.ErrNoCampaign
	}

	// validate available campaigns
	validCampaigns := []*models.Campaign{}
	copier.Copy(payloadValidator, m)
	copier.Copy(payloadValidator.Validators, m)

	for _, campaign := range campaigns {
		//  validate each campaign
		err = campaign.Validators.Validate(payloadValidator)

		if err == nil {
			validCampaigns = append(validCampaigns, campaign)
		}
	}

	if len(validCampaigns) < 1 {
		// no valid campaign available
		log.Error(err)

		return nil, models.ErrNoCampaign
	}

	// get latest campaign
	latestCampaign := validCampaigns[0]

	// get campaign formula
	expression, err := govaluate.NewEvaluableExpression(latestCampaign.Validators.Formula)

	if err != nil {
		log.Error(err)

		return nil, err
	}

	// get formula parameters
	parameters := make(map[string]interface{}, 8)
	parameters["transactionAmount"] = m.TransactionAmount
	parameters["multiplier"] = *latestCampaign.Validators.Multiplier
	parameters["value"] = *latestCampaign.Validators.Value
	result, err := expression.Evaluate(parameters)

	if err != nil {
		log.Error(err)

		return nil, err
	}

	// Parse interface to float
	parseFloat, err := getFloat(result)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	pointAmount := math.Floor(parseFloat)

	// store campaign transaction
	campaignTrx := &models.CampaignTrx{
		UserID:          m.UserID,
		PointAmount:     &pointAmount,
		TransactionType: models.TransactionPointTypeDebet,
		TransactionDate: &now,
		Campaign:        latestCampaign,
		CreatedAt:       &now,
	}

	err = cmpgn.campaignRepo.SavePoint(ctx, campaignTrx)

	if err != nil {
		return nil, err
	}

	p := new(models.UserPoint)
	p.UserPoint = pointAmount

	return p, nil
}

func (cmpgn *campaignUseCase) GetUserPoint(c context.Context, userID string) (*models.UserPoint, error) {
	ctx, cancel := context.WithTimeout(c, cmpgn.contextTimeout)
	defer cancel()
	pointAmount, err := cmpgn.campaignRepo.GetUserPoint(ctx, userID)

	if err != nil {
		return nil, err
	}

	p := new(models.UserPoint)
	p.UserPoint = pointAmount

	return p, nil
}

func (cmpgn *campaignUseCase) GetUserPointHistory(c context.Context, userID string) ([]models.CampaignTrx, error) {
	ctx, cancel := context.WithTimeout(c, cmpgn.contextTimeout)
	defer cancel()
	dataHistory, err := cmpgn.campaignRepo.GetUserPointHistory(ctx, userID)

	if err != nil {
		return nil, err
	}

	return dataHistory, nil
}

func getFloat(unk interface{}) (float64, error) {
	v := reflect.ValueOf(unk)
	v = reflect.Indirect(v)

	if !v.Type().ConvertibleTo(floatType) {
		return 0, fmt.Errorf("cannot convert %v to float64", v.Type())
	}

	fv := v.Convert(floatType)
	return fv.Float(), nil
}
