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

	return strconv.Itoa(countCampaign), listCampaign, nil
}

func (cmpgn *campaignUseCase) GetCampaignValue(c context.Context, m *models.GetCampaignValue) (*models.UserPoint, error) {
	ctx, cancel := context.WithTimeout(c, cmpgn.contextTimeout)
	defer cancel()

	dataCampaign, err := cmpgn.campaignRepo.GetValidatorCampaign(ctx, m)
	if err != nil {
		return nil, models.ErrNoCampaign
	}

	// Calculate point
	expression, err := govaluate.NewEvaluableExpression(dataCampaign.Validators.Formula)
	parameters := make(map[string]interface{}, 8)
	parameters["transactionAmount"] = m.TransactionAmount
	parameters["multiplier"] = dataCampaign.Validators.Multiplier
	parameters["value"] = dataCampaign.Validators.Value
	result, err := expression.Evaluate(parameters)

	// Parse interface to float
	parseFloat, err := getFloat(result)
	pointAmount := math.Floor(parseFloat)

	campaignTrx := &models.CampaignTrx{
		UserID:          m.UserId,
		PointAmount:     &pointAmount,
		TransactionType: models.TransactionPointTypeDebet,
		TransactionDate: &models.TimeNow,
		Campaign:        dataCampaign,
		CreatedAt:       &models.TimeNow,
	}

	err = cmpgn.campaignRepo.SavePoint(ctx, campaignTrx)
	if err != nil {
		return nil, err
	}

	p := new(models.UserPoint)
	p.UserPoint = pointAmount

	return p, nil
}

func (cmpgn *campaignUseCase) GetUserPoint(c context.Context, userId string) (*models.UserPoint, error) {
	ctx, cancel := context.WithTimeout(c, cmpgn.contextTimeout)
	defer cancel()

	pointAmount, err := cmpgn.campaignRepo.GetUserPoint(ctx, userId)
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
