package usecase

import (
	"context"
	"errors"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo"
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

func (cmpgn *campaignUseCase) CreateCampaign(c echo.Context, campaign *models.Campaign) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	err := cmpgn.campaignRepo.CreateCampaign(c, campaign)

	if err != nil {
		requestLogger.Debug(models.ErrCampaignFailed)

		return models.ErrCampaignFailed
	}

	return nil
}

func (cmpgn *campaignUseCase) UpdateCampaign(c echo.Context, id string, updateCampaign *models.UpdateCampaign) error {
	var campaignDetail *models.Campaign
	now := time.Now()
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	campaignID, err := strconv.Atoi(id)

	if err != nil {
		requestLogger.Debug(err)

		return errors.New("Something went wrong with input ID")
	}

	campaignDetail, err = cmpgn.campaignRepo.GetCampaignDetail(c, int64(campaignID))

	if err != nil {
		requestLogger.Debug(models.ErrNoCampaign)

		return models.ErrNoCampaign
	}

	vEndDate, _ := time.Parse(time.RFC3339, campaignDetail.EndDate)

	if vEndDate.Before(now.Add(time.Hour * -24)) {
		requestLogger.Debug(models.ErrCampaignExpired)

		return models.ErrCampaignExpired
	}

	err = cmpgn.campaignRepo.UpdateCampaign(c, int64(campaignID), updateCampaign)

	if err != nil {
		requestLogger.Debug(models.ErrCampaignUpdateFailed)

		return models.ErrCampaignUpdateFailed
	}

	return nil
}

func (cmpgn *campaignUseCase) GetCampaign(c context.Context, name string, status string, startDate string, endDate string, page int, limit int) (string, []*models.Campaign, error) {
	ctx, cancel := context.WithTimeout(c, cmpgn.contextTimeout)
	defer cancel()

	err := cmpgn.campaignRepo.UpdateExpiryDate(ctx)

	if err != nil {
		log.Debug("Update Status Base on Expiry Date: ", err)
	}

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

func (cmpgn *campaignUseCase) GetCampaignDetail(c echo.Context, id string) (*models.Campaign, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	campaignID, err := strconv.Atoi(id)

	if err != nil {
		requestLogger.Debug(err)

		return nil, errors.New("Something went wrong with input ID")
	}

	campaignDetail, err := cmpgn.campaignRepo.GetCampaignDetail(c, int64(campaignID))

	if err != nil {
		requestLogger.Debug(models.ErrNoCampaign)

		return nil, models.ErrNoCampaign
	}

	return campaignDetail, nil
}

func (cmpgn *campaignUseCase) GetCampaignValue(c context.Context, m *models.GetCampaignValue) (*models.UserPoint, error) {
	now := time.Now()
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
		UserID:          m.UserID,
		PointAmount:     &pointAmount,
		TransactionType: models.TransactionPointTypeDebet,
		TransactionDate: &now,
		Campaign:        dataCampaign,
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

func (cmpgn *campaignUseCase) UpdateStatusBasedOnStartDate() error {

	err := cmpgn.campaignRepo.UpdateStatusBasedOnStartDate()
	if err != nil {
		log.Debug("Update Status Base on Start Date: ", err)
		return err
	}
	return nil
}
