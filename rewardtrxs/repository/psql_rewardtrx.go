package repository

import (
	"database/sql"
	"encoding/json"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewardtrxs"
	"math/rand"
	"time"

	"github.com/labstack/echo"
)

type psqlRewardTrxRepository struct {
	Conn *sql.DB
}

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// NewPsqlRewardTrxRepository will create an object that represent the rewardtrxs.Repository interface
func NewPsqlRewardTrxRepository(Conn *sql.DB) rewardtrxs.Repository {
	return &psqlRewardTrxRepository{Conn}
}

func (rwdTrxRepo *psqlRewardTrxRepository) Create(c echo.Context, payload models.PayloadValidator, rewardID int64, resp []models.RewardResponse) (models.RewardTrx, error) {
	var rewardTrx models.RewardTrx
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	refID := randRefID(20)

	query := `INSERT INTO reward_transactions (status, ref_id, cif, reward_id, used_promo_code, transaction_date, inquired_date, request_data, response_data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	stmt, err := rwdTrxRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return rewardTrx, err
	}

	var lastID int64
	requestData, err := json.Marshal(payload)
	responseData, err := json.Marshal(models.RewardsInquiry{RefTrx: refID, Rewards: &resp})

	if err != nil {
		requestLogger.Debug(err)

		return rewardTrx, err
	}

	trxDate, err := time.Parse(time.RFC3339, payload.TransactionDate)

	if err != nil {
		requestLogger.Debug(models.ErrTrxDateFormat)

		return rewardTrx, err
	}

	rewardTrx = models.RewardTrx{
		Status:          &models.RewardTrxInquired,
		RefID:           refID,
		RewardID:        &rewardID,
		CIF:             payload.CIF,
		UsedPromoCode:   payload.PromoCode,
		TransactionDate: &trxDate,
		InquiryDate:     &now,
		RequestData:     string(requestData),
		ResponseData:    string(responseData),
		CreatedAt:       &now,
	}

	err = stmt.QueryRow(
		&rewardTrx.Status, &rewardTrx.RefID, &rewardTrx.CIF, &rewardTrx.RewardID, &rewardTrx.UsedPromoCode, &rewardTrx.TransactionDate,
		&rewardTrx.InquiryDate, &rewardTrx.RequestData, &rewardTrx.ResponseData, &rewardTrx.CreatedAt,
	).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return rewardTrx, err
	}

	return rewardTrx, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) GetByRefID(c echo.Context, refID string) (models.RewardsInquiry, error) {
	var rewardInquiry models.RewardsInquiry
	var rwdInquiry json.RawMessage
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT response_data from reward_transactions where status = $1 and ref_id = $2`
	err := rwdTrxRepo.Conn.QueryRow(query, models.RewardTrxInquired, refID).Scan(&rwdInquiry)

	if err != nil {
		requestLogger.Debug(err)

		return rewardInquiry, err
	}

	err = json.Unmarshal([]byte(rwdInquiry), &rewardInquiry)

	if err != nil {
		requestLogger.Debug(err)

		return rewardInquiry, err
	}

	return rewardInquiry, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) UpdateSuccess(c echo.Context, payload map[string]interface{}) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	status := int8(1)
	refCore := payload["refCore"].(string)
	cif := payload["cif"].(string)
	refID := payload["refTrx"].(string)
	query := `UPDATE reward_transactions SET status = $1, ref_core = $2, successed_date = $3, updated_date = $4 WHERE cif = $5 and ref_id = $6`
	stmt, err := rwdTrxRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	var lastID int64

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(&status, &refCore, &now, &now, &cif, &refID).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) UpdateReject(c echo.Context, payload map[string]interface{}) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	status := int8(2)
	cif := payload["cif"].(string)
	refID := payload["refTrx"].(string)
	query := `UPDATE reward_transactions SET status = $1, rejected_date = $2, updated_date = $3 WHERE cif = $4 and ref_id = $5`
	stmt, err := rwdTrxRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	var lastID int64

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(&status, &now, &now, &cif, &refID).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func randRefID(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
