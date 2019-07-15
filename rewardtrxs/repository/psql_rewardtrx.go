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

// NewPsqlRewardTrxRepository will create an object that represent the quotas.Repository interface
func NewPsqlRewardTrxRepository(Conn *sql.DB) rewardtrxs.Repository {
	return &psqlRewardTrxRepository{Conn}
}

func (quotTrxRepo *psqlRewardTrxRepository) Create(c echo.Context, payload map[string]interface{}, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	status := int8(0)
	refID := RandRefID(20)
	query := `INSERT INTO reward_transactions (status, ref_id, cif, reward_id, used_promo_code, inquiry_date, timeout_date, request_data, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	stmt, err := quotTrxRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	var lastID int64
	requestData, err := json.Marshal(payload)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(&status, &refID, payload["cif"], &rewardID, payload["promoCode"], &now, &now, string(requestData), &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

const letterBytes = "1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandRefID for unique refID based on time
func RandRefID(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}

func (quotTrxRepo *psqlRewardTrxRepository) UpdateSuccess(c echo.Context, payload map[string]interface{}) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	status := int8(1)
	refCore := payload["refCore"].(string)
	cif := payload["cif"].(string)
	refID := payload["refTrx"].(string)
	query := `UPDATE reward_transactions SET status = $1, ref_core = $2, successed_date = $3, updated_date = $4 WHERE cif = $5 and ref_id = $6`
	stmt, err := quotTrxRepo.Conn.Prepare(query)

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

func (quotTrxRepo *psqlRewardTrxRepository) UpdateReject(c echo.Context, payload map[string]interface{}) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	status := int8(2)
	cif := payload["cif"].(string)
	refID := payload["refTrx"].(string)
	query := `UPDATE reward_transactions SET status = $1, rejected_date = $2, updated_date = $3 WHERE cif = $4 and ref_id = $5`
	stmt, err := quotTrxRepo.Conn.Prepare(query)

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
