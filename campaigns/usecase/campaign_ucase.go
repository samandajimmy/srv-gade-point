package usecase

import (
	"context"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"math"
	"reflect"
	"time"

	govaluate "gopkg.in/Knetic/govaluate.v2"
)

var floatType = reflect.TypeOf(float64(0))

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

func (a *campaignUseCase) GetCampaignValue(c context.Context, m *models.GetCampaignValue) (*models.UserPoint, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	dataCampaign, err := a.campaignRepo.GetValidatorCampaign(ctx, m)
	if err != nil {
		return nil, models.ErrNoCampaign
	}

	//Calculate point
	expression, err := govaluate.NewEvaluableExpression(dataCampaign.Validators.Formula)
	parameters := make(map[string]interface{}, 8)
	parameters["transactionAmount"] = m.TransactionAmount
	parameters["multiplier"] = dataCampaign.Validators.Multiplier
	parameters["value"] = dataCampaign.Validators.Value
	result, err := expression.Evaluate(parameters)

	//Parse interface to float
	parseFloat, err := getFloat(result)
	pointAmount := math.Floor(parseFloat)

	saveTransactionPoint := &models.SaveTransactionPoint{
		UserId:          m.UserId,
		PointAmount:     pointAmount,
		TransactionType: models.TransactionPointTypeDebet,
		TransactionDate: time.Now(),
		CampaingId:      dataCampaign.ID,
		PromoCodeId:     0,
		CreatedAt:       time.Now(),
	}

	err = a.campaignRepo.SavePoint(ctx, saveTransactionPoint)
	if err != nil {
		return nil, err
	}

	p := new(models.UserPoint)
	p.UserPoint = pointAmount

	return p, nil
}

func (a *campaignUseCase) GetUserPoint(c context.Context, userId string) (*models.UserPoint, error) {
	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	pointAmount, err := a.campaignRepo.GetUserPoint(ctx, userId)
	if err != nil {
		return nil, err
	}

	p := new(models.UserPoint)
	p.UserPoint = pointAmount

	return p, nil
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
