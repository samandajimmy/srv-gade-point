package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewardtrxs"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type psqlRewardTrxRepository struct {
	Conn *sql.DB
}

// NewPsqlRewardTrxRepository will create an object that represent the rewardtrxs.Repository interface
func NewPsqlRewardTrxRepository(Conn *sql.DB) rewardtrxs.RtRepository {
	return &psqlRewardTrxRepository{Conn}
}

func (rwdTrxRepo *psqlRewardTrxRepository) Create(c echo.Context, payload models.PayloadValidator,
	resp models.RewardsInquiry) ([]*models.RewardTrx, error) {

	var rewardTrx []*models.RewardTrx
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	respData := resp
	respData.Rewards = nil
	rootRefTrx := resp.RefTrx
	refID := rootRefTrx

	for _, rwdResp := range *resp.Rewards {
		var lastID int64
		var rwdRespData []models.RewardResponse
		var rewardID sql.NullInt64
		expTime, _ := strconv.ParseInt(os.Getenv(`REWARD_TRX_TIMEOUT`), 10, 64)
		timeoutDate := now.Add(time.Duration(expTime) * time.Minute)
		rwdRespData = append(rwdRespData, rwdResp)
		respData.Rewards = &rwdRespData

		query := `INSERT INTO reward_transactions (status, ref_id, cif, reward_id, used_promo_code,
			transaction_date, inquired_date, request_data, response_data, created_at, timeout_date,
			root_ref_id) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING id`

		stmt, err := rwdTrxRepo.Conn.Prepare(query)

		if err != nil {
			requestLogger.Debug(err)

			return rewardTrx, err
		}

		if rwdResp.RefTrx != "" {
			refID = rwdResp.RefTrx
		}

		requestData, _ := json.Marshal(payload)
		responseData, err := json.Marshal(respData)

		if err != nil {
			requestLogger.Debug(err)

			return rewardTrx, err
		}

		trxDate, err := time.Parse(models.DateTimeFormatMillisecond, payload.TransactionDate)

		if err != nil {
			requestLogger.Debug(models.ErrTrxDateFormat)

			return rewardTrx, err
		}

		if rwdResp.RewardID != 0 {
			_ = rewardID.Scan(rwdResp.RewardID)
		}

		rwdTrx := models.RewardTrx{
			Status:          &models.RewardTrxInquired,
			RefID:           refID,
			CIF:             payload.CIF,
			UsedPromoCode:   payload.PromoCode,
			TransactionDate: &trxDate,
			InquiredDate:    &now,
			CreatedAt:       &now,
			TimeoutDate:     &timeoutDate,
		}

		reqData := string(requestData)
		respData := string(responseData)

		err = stmt.QueryRow(
			&rwdTrx.Status, &rwdTrx.RefID, &rwdTrx.CIF, rewardID, &rwdTrx.UsedPromoCode,
			&rwdTrx.TransactionDate, &rwdTrx.InquiredDate, &reqData, &respData,
			&rwdTrx.CreatedAt, &rwdTrx.TimeoutDate, rootRefTrx,
		).Scan(&lastID)

		if err != nil {
			requestLogger.Debug(err)

			return rewardTrx, err
		}

		rewardTrx = append(rewardTrx, &rwdTrx)
	}

	return rewardTrx, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) GetByRefID(c echo.Context, refID string) (models.RewardsInquiry, error) {
	var rewardInquiry models.RewardsInquiry
	var rwdRespData []models.RewardResponse
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT response_data from reward_transactions where status = $1 and
		(ref_id = $2 or root_ref_id = $2)`
	rows, err := rwdTrxRepo.Conn.Query(query, models.RewardTrxInquired, refID)

	if err != nil {
		requestLogger.Debug(err)

		return rewardInquiry, err
	}

	defer rows.Close()

	for rows.Next() {
		var rwdInquiry json.RawMessage
		var rwdResp []models.RewardResponse

		err = rows.Scan(
			&rwdInquiry,
		)

		if err != nil {
			requestLogger.Debug(err)

			return rewardInquiry, err
		}

		err = json.Unmarshal([]byte(rwdInquiry), &rewardInquiry)

		if err != nil {
			requestLogger.Debug(err)

			return rewardInquiry, err
		}

		rwdResp = *rewardInquiry.Rewards
		rwdRespData = append(rwdRespData, rwdResp[0])
	}

	if (rewardInquiry != models.RewardsInquiry{}) {
		rewardInquiry.Rewards = &rwdRespData
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
	var rewardTrxReqData models.PayloadValidator
	var rewardTrxRespData models.RewardsInquiry
	var reward models.Reward
	var reqData json.RawMessage
	var respData json.RawMessage
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT rt.id, rt.status, coalesce(rt.ref_core, ''), rt.ref_id, rt.reward_id, rt.cif,
		rt.inquired_date, rt.transaction_date, rt.request_data, rt.response_data, r.type, rt.used_promo_code
		from reward_transactions rt
		left join rewards r on rt.reward_id = r.id where ref_id = $1`
	err := rwdTrxRepo.Conn.QueryRow(query, refID).Scan(
		&result.ID,
		&result.Status,
		&result.RefCore,
		&result.RefID,
		&result.RewardID,
		&result.CIF,
		&result.InquiredDate,
		&result.TransactionDate,
		&reqData,
		&respData,
		&reward.Type,
		&result.UsedPromoCode,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	_ = json.Unmarshal([]byte(reqData), &rewardTrxReqData)
	err = json.Unmarshal([]byte(respData), &rewardTrxRespData)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	result.Reward = &reward
	result.RequestData = &rewardTrxReqData
	result.ResponseData = &rewardTrxRespData

	return &result, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) CheckRootRefId(c echo.Context, rootRefID string) ([]string, error) {
	var result models.RewardTrx
	var results []string
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT rt.id,rt.ref_id
		from reward_transactions rt
		left join rewards r on rt.reward_id = r.id where rt.root_ref_id = $1`

	rows, err := rwdTrxRepo.Conn.Query(query, rootRefID)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(
			&result.ID,
			&result.RefID,
		)
		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		results = append(results, result.RefID)
	}

	return results, nil
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

func (rwdTrxRepo *psqlRewardTrxRepository) GetRewardByPayload(c echo.Context,
	payload models.PayloadValidator, voucherCode *models.VoucherCode) ([]*models.Reward, error) {

	var reward []*models.Reward
	var validator json.RawMessage
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `SELECT rt.ref_id, coalesce(r.id, 0) id, r.campaign_id, coalesce(r.journal_account, '') journal_account,
		r.type, coalesce(r.validators, '{}'::jsonb) validators, rt.root_ref_id
		FROM reward_transactions rt left join rewards r on rt.reward_id = r.id
		WHERE rt.status = $1 and rt.cif = $2 and rt.used_promo_code = $3 and rt.request_data->>'phone' = $4
		and rt.request_data->'validators'->>'product' = $5 and rt.request_data->'validators'->>'transactionType' = $6`

	rows, err := rwdTrxRepo.Conn.Query(query, models.RewardTrxInquired, payload.CIF, payload.PromoCode,
		payload.Phone, payload.Validators.Product, payload.Validators.TransactionType)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var r models.Reward

		err = rows.Scan(
			&r.RefID,
			&r.ID,
			&r.CampaignID,
			&r.JournalAccount,
			&r.Type,
			&validator,
			&r.RootRefID,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		err = json.Unmarshal([]byte(validator), &r.Validators)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		if voucherCode != nil {
			r.JournalAccount = voucherCode.Voucher.JournalAccount
			r.Type = voucherCode.Voucher.Type
			r.Validators = voucherCode.Voucher.Validators
		}

		reward = append(reward, &r)
	}

	return reward, nil
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

	// make voucher codes timeout
	rwdTrxRepo.voucherTrxTimeout(rewardTrx.RefID)

	if rewardTrx.RewardID == nil {
		return
	}

	rwdTrxRepo.updateLockedQuota(*rewardTrx.RewardID, rewardTrx.RefID)
}

func (rwdTrxRepo *psqlRewardTrxRepository) UpdateTimeoutTrx() error {
	now := time.Now()
	query := `UPDATE reward_transactions SET status = $1, updated_at = $2 where timeout_date <= $3
		and status = $4 RETURNING reward_id, ref_id`
	rows, err := rwdTrxRepo.Conn.Query(query, &models.RewardTrxTimeOut, &now, &now, &models.RewardTrxInquired)

	if err != nil {
		logrus.Debug(err)
	}

	defer rows.Close()

	for rows.Next() {
		var t models.RewardTrx
		err = rows.Scan(
			&t.RewardID,
			&t.RefID,
		)

		if err != nil {
			logrus.Debug(err)
		}

		// make voucher codes timeout
		rwdTrxRepo.voucherTrxTimeout(t.RefID)

		if t.RewardID != nil {
			rwdTrxRepo.updateLockedQuota(*t.RewardID, t.RefID)
		}
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
	query := `UPDATE voucher_codes SET status = $1, user_id = $2, bought_date = NULL,
		updated_at = $3 where status = $4 and ref_id = $5`
	stmt, err := rwdTrxRepo.Conn.Prepare(query)

	if err != nil {
		logrus.Debug(err)
	}

	rowVC, err := stmt.Query(&models.VoucherCodeStatusAvailable, "", &now,
		models.VoucherCodeStatusBooked, &refID)

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

func (rwdTrxRepo *psqlRewardTrxRepository) CountRewardTrxs(c echo.Context, rewardPayload *models.RewardsPayload) (int64, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var where string
	var count int64

	query := `SELECT count(ID) FROM reward_transactions`

	if rewardPayload.RewardID != "" {
		where += " WHERE reward_id = '" + rewardPayload.RewardID + "'"
	}

	if rewardPayload.StartTransactionDate != "" {
		where += " AND transaction_date::timestamp::date >= '" + rewardPayload.StartTransactionDate + "'"
	}

	if rewardPayload.EndTransactionDate != "" {
		where += " AND transaction_date::timestamp::date <= '" + rewardPayload.EndTransactionDate + "'"
	}

	if rewardPayload.StartSuccededDate != "" {
		where += " AND succeeded_date::timestamp::date >= '" + rewardPayload.StartSuccededDate + "'"
	}

	if rewardPayload.EndTransactionDate != "" {
		where += " AND succeeded_date::timestamp::date <= '" + rewardPayload.EndTransactionDate + "'"
	}

	if rewardPayload.Status != "" {
		where += " AND status = '" + rewardPayload.Status + "'"
	}

	query += where + ";"
	err := rwdTrxRepo.Conn.QueryRow(query).Scan(&count)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	return count, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) GetRewardTrxs(c echo.Context, rewardPayload *models.RewardsPayload) ([]models.RewardTrx, error) {
	var data []models.RewardTrx
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var paging, where string

	query := `SELECT id, status, coalesce(ref_core, ''), ref_id, reward_id, cif, used_promo_code, inquired_date, succeeded_date, rejected_date,
		timeout_date, transaction_date FROM reward_transactions`

	if rewardPayload.RewardID != "" {
		where += " WHERE reward_id = '" + rewardPayload.RewardID + "'"
	}

	if rewardPayload.StartTransactionDate != "" {
		where += " AND transaction_date::timestamp::date >= '" + rewardPayload.StartTransactionDate + "'"
	}

	if rewardPayload.EndTransactionDate != "" {
		where += " AND transaction_date::timestamp::date <= '" + rewardPayload.EndTransactionDate + "'"
	}

	if rewardPayload.StartSuccededDate != "" {
		where += " AND succeeded_date::timestamp::date >= '" + rewardPayload.StartSuccededDate + "'"
	}

	if rewardPayload.EndSuccededDate != "" {
		where += " AND succeeded_date::timestamp::date <= '" + rewardPayload.EndSuccededDate + "'"
	}

	if rewardPayload.Status != "" {
		where += " AND status = '" + rewardPayload.Status + "'"
	}

	if rewardPayload.Page > 0 || rewardPayload.Limit > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", rewardPayload.Limit, (rewardPayload.Page-1)*rewardPayload.Limit)
	}

	query += where + " order by transaction_date desc" + paging + ";"
	rows, err := rwdTrxRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var rwdTrx models.RewardTrx
		var inquiredDate, succeededDate, rejectedDate, timeoutDate pq.NullTime

		err = rows.Scan(
			&rwdTrx.ID,
			&rwdTrx.Status,
			&rwdTrx.RefCore,
			&rwdTrx.RefID,
			&rwdTrx.RewardID,
			&rwdTrx.CIF,
			&rwdTrx.UsedPromoCode,
			&inquiredDate,
			&succeededDate,
			&rejectedDate,
			&timeoutDate,
			&rwdTrx.TransactionDate,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		rwdTrx.InquiredDate = &inquiredDate.Time
		rwdTrx.SucceededDate = &succeededDate.Time
		rwdTrx.RejectedDate = &rejectedDate.Time
		rwdTrx.TimeoutDate = &timeoutDate.Time

		data = append(data, rwdTrx)
	}

	return data, nil
}

func (rwdTrxRepo *psqlRewardTrxRepository) voucherTrxTimeout(refId string) {
	now := time.Now()
	query := `UPDATE voucher_codes SET user_id = NULL, ref_id = NULL, status = $1, updated_at = $2
		where status = $3 and ref_id = $4`
	stmt, err := rwdTrxRepo.Conn.Prepare(query)

	if err != nil {
		logrus.Debug(err)
	}

	rows, err := stmt.Query(&models.VoucherCodeStatusAvailable, &now, &models.VoucherCodeStatusInquired, refId)

	if err != nil {
		logrus.Debug(err)
	}

	defer rows.Close()
}
