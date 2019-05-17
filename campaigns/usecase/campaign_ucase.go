package usecase

import (
	"errors"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"math"
	"strconv"
	"time"

	"github.com/jinzhu/copier"
	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/sirupsen/logrus"
)

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
	var result float64
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	payloadValidator := &models.PayloadValidator{}
	payloadValidator.Validators = &models.Validator{}
	now := time.Now()
	dataCampaign, err := cmpgn.campaignRepo.GetValidatorCampaign(c, payload)

	// get available campaign
	campaigns, err := cmpgn.campaignRepo.GetCampaignAvailable(c)

	if err != nil {
		requestLogger.Debug(models.ErrNoCampaign)

		return nil, models.ErrNoCampaign
	}

	// validate available campaigns
	validCampaigns := []*models.Campaign{}
	copier.Copy(payloadValidator, payload)
	copier.Copy(payloadValidator.Validators, payload)

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
	if payloadValidator.Validators.Formula == "" {
		result = float64(0)
	} else {
		result, err = latestCampaign.Validators.GetFormulaResult(payloadValidator)
	}

	if err != nil {
		log.Error(err)
		return nil, err
	}

	pointAmount := math.Floor(result)

	// store campaign transaction
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
	p.UserPoint = &pointAmount

	return p, nil
}

func (cmpgn *campaignUseCase) GetUserPoint(c echo.Context, userID string) (*models.UserPoint, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	p := new(models.UserPoint)
	zero := float64(0)
	pointAmount, err := cmpgn.campaignRepo.GetUserPoint(c, userID)

	if err != nil {
		requestLogger.Debug(models.ErrGetUserPoint)
		p.UserPoint = &zero

		return p, models.ErrGetUserPoint
	}

	if pointAmount == 0 {
		requestLogger.Debug(models.ErrUserPointNA)
		p.UserPoint = &zero

		return p, models.ErrUserPointNA
	}

	p.UserPoint = &pointAmount

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

func (cmpgn *campaignUseCase) UpdateStatusBasedOnStartDate() error {
	err := cmpgn.campaignRepo.UpdateStatusBasedOnStartDate()

	if err != nil {
		logrus.Debug("Update Status Base on Start Date: ", err)

		return err
	}
	return nil
}
