package repository

import (
	"database/sql"
	"encoding/json"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewardtrxs"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
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
	expTime, _ := strconv.ParseInt(os.Getenv(`REWARD_TRX_TIMEOUT`), 10, 64)
	timeoutDate := now.Add(time.Duration(expTime) * time.Minute)

	query := `INSERT INTO reward_transactions (status, ref_id, cif, reward_id, used_promo_code, transaction_date, inquired_date, request_data, response_data, created_at, timeout_date)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`

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

	trxDate, err := time.Parse(models.DateTimeFormatMillisecond, payload.TransactionDate)

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
		InquiredDate:    &now,
		RequestData:     string(requestData),
		ResponseData:    string(responseData),
		CreatedAt:       &now,
		TimeoutDate:     &timeoutDate,
	}

	err = stmt.QueryRow(
		&rewardTrx.Status, &rewardTrx.RefID, &rewardTrx.CIF, &rewardTrx.RewardID, &rewardTrx.UsedPromoCode, &rewardTrx.TransactionDate,
		&rewardTrx.InquiredDate, &rewardTrx.RequestData, &rewardTrx.ResponseData, &rewardTrx.CreatedAt, &rewardTrx.TimeoutDate,
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

func (rwdTrxRepo *psqlRewardTrxRepository) CheckTrx(c echo.Context, refID string) (*models.RewardTrx, error) {
	var result models.RewardTrx
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT id, status, coalesce(ref_core, ''), ref_id, reward_id, cif, inquired_date, transaction_date from reward_transactions where ref_id = $1 and status = $2`
	err := rwdTrxRepo.Conn.QueryRow(query, refID, models.RewardTrxInquired).Scan(
		&result.ID,
		&result.Status,
		&result.RefCore,
		&result.RefID,
		&result.RewardID,
		&result.CIF,
		&result.InquiredDate,
		&result.TransactionDate,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return &result, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) CheckRefID(c echo.Context, refID string) (*models.RewardTrx, error) {
	var result models.RewardTrx
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT id, status, coalesce(ref_core, ''), ref_id, reward_id, cif, inquired_date, transaction_date from reward_transactions where ref_id = $1`
	err := rwdTrxRepo.Conn.QueryRow(query, refID).Scan(
		&result.ID,
		&result.Status,
		&result.RefCore,
		&result.RefID,
		&result.RewardID,
		&result.CIF,
		&result.InquiredDate,
		&result.TransactionDate,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return &result, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) UpdateRewardTrx(c echo.Context, rwdPayment *models.RewardPayment, status int64) error {
	var refCore, cif string
	var succeedDate, rejectedDate *time.Time
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()

	if status == models.RewardTrxSucceeded {
		refCore = rwdPayment.RefCore
		cif = rwdPayment.CIF
		succeedDate = &now
	} else if status == models.RewardTrxTimeOutForceToSucceeded {
		refCore = rwdPayment.RefCore
		cif = rwdPayment.CIF
		succeedDate = &now
	} else {
		rejectedDate = &now
		cif = rwdPayment.CIF
	}

	query := `UPDATE reward_transactions SET cif = $1, ref_core = $2 , status = $3, succeeded_date = $4, rejected_date = $5,
		updated_at = $6 where ref_id = $7`

	stmt, err := rwdTrxRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	rows, err := stmt.Query(&cif, &refCore, &status, succeedDate, rejectedDate, &now, &rwdPayment.RefTrx)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	defer rows.Close()

	return nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) CountByCIF(c echo.Context, quot models.Quota, rwd models.Reward, cif string) (int64, error) {
	var startDate, endDate time.Time
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	if quot.NumberOfDays != nil && *quot.NumberOfDays == models.CampaignPeriod {
		startDate, _ = time.Parse(time.RFC3339, rwd.Campaign.StartDate)
		endDate, _ = time.Parse(time.RFC3339, rwd.Campaign.EndDate)
	} else {
		startDate = *quot.LastCheck
		endDate = quot.LastCheck.AddDate(0, 0, int(*quot.NumberOfDays-1))
	}

	query := `select count(ID) from reward_transactions where cif = $1 and transaction_date::date >= $2 and transaction_date::date <= $3 and status in($4, $5, $6) and reward_id = $7`
	stmt, err := rwdTrxRepo.Conn.Prepare(query)
	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	var counter int64

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	err = stmt.QueryRow(&cif, startDate, endDate, models.RewardTrxInquired, models.RewardTrxSucceeded, models.RewardTrxTimeOutForceToSucceeded, rwd.ID).Scan(&counter)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	return counter, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) GetRewardByPayload(c echo.Context, payload models.PayloadValidator) (*models.Reward, string, error) {
	result := new(models.Reward)
	refID := ""
	var validator json.RawMessage
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `SELECT r.id, r.campaign_id, r.journal_account, r.type, r.validators, rt.ref_id
		FROM reward_transactions rt  join rewards r on rt.reward_id = r.id WHERE rt.status = $1 and rt.cif= $2 and rt.used_promo_code= $3`

	err := rwdTrxRepo.Conn.QueryRow(query, models.RewardTrxInquired, payload.CIF, payload.PromoCode).Scan(
		&result.ID,
		&result.CampaignID,
		&result.JournalAccount,
		&result.Type,
		&validator,
		&refID,
	)

	if refID != "" {
		err = json.Unmarshal([]byte(validator), &result.Validators)
		if err != nil {
			requestLogger.Debug(err)

			return nil, "", err
		}
	}

	if err != nil {
		requestLogger.Debug(err)

		return nil, "", err
	}

	return result, refID, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) RewardTrxTimeout(rewardTrx models.RewardTrx) {
	now := time.Now()
	query := `UPDATE reward_transactions SET status = $1, updated_at = $2 where status = $3 and ref_id = $4`
	stmt, err := rwdTrxRepo.Conn.Prepare(query)

	if err != nil {
		logrus.Debug(err)
	}

	rows, err := stmt.Query(&models.RewardTrxTimeOut, &now, &models.RewardTrxInquired, &rewardTrx.RefID)

	if err != nil {
		logrus.Debug(err)
	}

	defer rows.Close()

	rwdTrxRepo.updateLockedQuota(*rewardTrx.RewardID, rewardTrx.RefID)
}

func (rwdTrxRepo *psqlRewardTrxRepository) UpdateTimeoutTrx() error {
	now := time.Now()
	query := `UPDATE reward_transactions SET status = $1, updated_at = $2 where timeout_date <= $3 and status = $4`
	rows, err := rwdTrxRepo.Conn.Query(query, &models.RewardTrxTimeOut, &now, &now, &models.RewardTrxInquired)

	if err != nil {
		logrus.Debug(err)
	}

	defer rows.Close()

	for rows.Next() {
		var t models.RewardTrx
		err = rows.Scan(
			&t.RewardID,
			&t.RefCore,
		)

		if err != nil {
			logrus.Debug(err)
		}

		rwdTrxRepo.updateLockedQuota(*t.RewardID, t.RefID)
	}

	if err != nil {
		logrus.Debug(err)
	}

	return nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) GetInquiredTrx() ([]models.RewardTrx, error) {
	var result []models.RewardTrx
	now := time.Now()
	query := `select id, status, coalesce(ref_core, ''), ref_id, reward_id, cif, used_promo_code, inquired_date, succeeded_date, timeout_date,
		rejected_date, transaction_date from reward_transactions
		where timeout_date >= $1 and status = $2`

	rows, err := rwdTrxRepo.Conn.Query(query, &now, &models.RewardTrxInquired)

	if err != nil {
		logrus.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t models.RewardTrx
		err = rows.Scan(
			&t.ID,
			&t.Status,
			&t.RefCore,
			&t.RefID,
			&t.RewardID,
			&t.CIF,
			&t.UsedPromoCode,
			&t.InquiredDate,
			&t.SucceededDate,
			&t.TimeoutDate,
			&t.RejectedDate,
			&t.TransactionDate,
		)

		if err != nil {
			logrus.Debug(err)

			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) updateLockedQuota(rewardID int64, refID string) {
	now := time.Now()
	zero := int64(0)
	query := `UPDATE voucher_codes SET status = $1, user_id = $2, bought_date = NULL, updated_at = $3 where status = 1 and ref_id = $4`
	stmt, err := rwdTrxRepo.Conn.Prepare(query)

	if err != nil {
		logrus.Debug(err)
	}

	rowVC, err := stmt.Query(&zero, "", &now, &refID)

	if err != nil {
		logrus.Debug(err)
	}

	defer rowVC.Close()

	query = `UPDATE quotas SET available = available + 1 WHERE reward_id = $1 and is_per_user = $2`
	stmt, err = rwdTrxRepo.Conn.Prepare(query)

	if err != nil {
		logrus.Debug(err)
	}

	rowQ, err := stmt.Query(rewardID, models.IsPerUserFalse)

	if err != nil {
		logrus.Debug(err)
	}

	defer rowQ.Close()
}

func randRefID(n int) string {
	rand.Seed(time.Now().UnixNano())

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}
	return string(b)
}
