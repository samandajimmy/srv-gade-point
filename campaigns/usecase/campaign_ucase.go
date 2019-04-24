package usecase

import (
	"errors"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"math"
	"reflect"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
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

func (cmpgn *campaignUseCase) GetCampaign(c echo.Context, payload map[string]interface{}) (string, []*models.Campaign, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	page, err := strconv.Atoi(payload["page"].(string))

	if err != nil {
		requestLogger.Debug(err)

		return "", nil, errors.New("Something went wrong with input page")
	}

	limit, err := strconv.Atoi(payload["limit"].(string))

	if err != nil {
		requestLogger.Debug(err)

		return "", nil, errors.New("Something went wrong with input limit")
	}

	payload["page"] = page
	payload["limit"] = limit
	_ = cmpgn.campaignRepo.UpdateExpiryDate(c)
	listCampaign, err := cmpgn.campaignRepo.GetCampaign(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetCampaign)

		return "", nil, models.ErrGetCampaign
	}

	countCampaign, err := cmpgn.campaignRepo.CountCampaign(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetCampaign)

		return "", nil, err
	}

	if countCampaign <= 0 {
		requestLogger.Debug(models.ErrGetCampaignCounter)

		return "", listCampaign, models.ErrGetCampaignCounter
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

func (cmpgn *campaignUseCase) GetCampaignValue(c echo.Context, payload *models.GetCampaignValue) (*models.UserPoint, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	dataCampaign, err := cmpgn.campaignRepo.GetValidatorCampaign(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrNoCampaign)

		return nil, models.ErrNoCampaign
	}

	// Calculate point
	expression, err := govaluate.NewEvaluableExpression(dataCampaign.Validators.Formula)

	if err != nil {
		requestLogger.Debug(models.ErrCalculateFormulaCampaign)

		return nil, models.ErrCalculateFormulaCampaign
	}

	parameters := make(map[string]interface{}, 8)
	parameters["transactionAmount"] = payload.TransactionAmount
	parameters["multiplier"] = dataCampaign.Validators.Multiplier
	parameters["value"] = dataCampaign.Validators.Value
	result, err := expression.Evaluate(parameters)

	if err != nil {
		requestLogger.Debug(err)

		return nil, models.ErrCalculateFormulaCampaign
	}

	// Parse interface to float
	parseFloat, err := getFloat(result)

	if err != nil {
		requestLogger.Debug(err)

		return nil, models.ErrCalculateFormulaCampaign
	}

	if math.IsInf(parseFloat, 0) {
		requestLogger.Debug("the result of formula is infinity and beyond")

		return nil, models.ErrCalculateFormulaCampaign
	}

	pointAmount := math.Floor(parseFloat)

	campaignTrx := &models.CampaignTrx{
		UserID:          payload.UserID,
		PointAmount:     &pointAmount,
		TransactionType: models.TransactionPointTypeDebet,
		TransactionDate: &now,
		ReffCore:        payload.ReffCore,
		Campaign:        dataCampaign,
		CreatedAt:       &now,
	}

	err = cmpgn.campaignRepo.SavePoint(c, campaignTrx)

	if err != nil {
		requestLogger.Debug(models.ErrStoreCampaignTrx)

		return nil, models.ErrStoreCampaignTrx
	}

	p := new(models.UserPoint)
	p.UserPoint = pointAmount

	return p, nil
}

func (cmpgn *campaignUseCase) GetUserPoint(c echo.Context, userID string) (*models.UserPoint, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	pointAmount, err := cmpgn.campaignRepo.GetUserPoint(c, userID)

	if err != nil {
		requestLogger.Debug(models.ErrGetUserPoint)

		return nil, models.ErrGetUserPoint
	}

	if pointAmount == 0 {
		requestLogger.Debug(models.ErrUserPointNA)

		return nil, models.ErrUserPointNA
	}

	p := new(models.UserPoint)
	p.UserPoint = pointAmount

	return p, nil
}

func (cmpgn *campaignUseCase) GetUserPointHistory(c echo.Context, payload map[string]interface{}) ([]models.CampaignTrx, string, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	counter, err := cmpgn.campaignRepo.CountUserPointHistory(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrUserPointHistoryNA)

		return nil, "", models.ErrUserPointHistoryNA
	}

	dataHistory, err := cmpgn.campaignRepo.GetUserPointHistory(c, payload)

	if err != nil {
		requestLogger.Debug(models.ErrGetUserPointHistory)

		return nil, "", models.ErrGetUserPointHistory
	}

	return dataHistory, counter, nil
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
		logrus.Debug("Update Status Base on Start Date: ", err)

		return err
	}
	return nil
}
