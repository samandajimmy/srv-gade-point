package repository

import (
	"database/sql"
	gcdb "gade/srv-gade-point/database"
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"time"

	"github.com/labstack/echo"
)

type psqlReferralsRepository struct {
	Conn *sql.DB
	Bun  *gcdb.DbBun
}

// NewPsqlReferralRepository will create an object that represent the referrals.RefRepository interface
func NewPsqlReferralRepository(Conn *sql.DB, Bun *gcdb.DbBun) referrals.RefRepository {
	return &psqlReferralsRepository{Conn, Bun}
}

func (refRepo *psqlReferralsRepository) RPostCoreTrx(c echo.Context, coreTrx []models.CoreTrxPayload) error {
	var nilFilters []string
	createdAt := time.Now()
	trxType := 1
	totalReward := 0
	stmts := []*gcdb.PipelineStmt{}
	for _, trx := range coreTrx {

		stmts = append(stmts, gcdb.NewPipelineStmt(`INSERT INTO core_transactions 
		(created_at, transaction_amount, loan_amount, interest_amount, product_code, 
		transaction_date, total_reward, transaction_id, marketing_code, transaction_type) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
			nilFilters, createdAt, trx.TrxAmount, trx.LoanAmount,
			trx.InterestAmount, trx.ProductCode, trx.TrxDate, totalReward, trx.TrxID,
			trx.MarketingCode, trxType))
	}

	err := gcdb.WithTransaction(refRepo.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	return nil
}

func (m *psqlReferralsRepository) RCreateReferral(c echo.Context, refcodes models.ReferralCodes) (models.ReferralCodes, error) {
	now := time.Now()
	query := `INSERT INTO referral_codes (cif, referral_code, campaign_id, created_at) 
		VALUES (?0, ?1, ?2, ?3) RETURNING id`

	_, err := m.Bun.QueryContext(c.Request().Context(), query, refcodes.CIF, refcodes.ReferralCode, refcodes.CampaignId, now)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ReferralCodes{}, err
	}

	refcodes.CreatedAt = now

	return refcodes, nil
}

func (m *psqlReferralsRepository) RGetReferralByCif(c echo.Context, cif string) (models.ReferralCodes, error) {

	var result models.ReferralCodes

	query := `select cif, referral_code, campaign_id, created_at, updated_at from referral_codes rc 
		where cif = ? order by created_at desc limit 1;`

	rows, err := m.Bun.QueryContext(c.Request().Context(), query, cif)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ReferralCodes{}, err
	}

	err = m.Bun.ScanRows(c.Request().Context(), rows, &result)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ReferralCodes{}, err
	}

	return result, nil
}

func (m *psqlReferralsRepository) RGetCampaignId(c echo.Context, prefix string) (int64, error) {

	var code int64
	now := time.Now()

	query := `select id from campaigns c where metadata->>'prefix' = ?0 and start_date <= ?1 
		and end_date >= ?1 order by created_at desc limit 1;`

	rows, err := m.Bun.QueryContext(c.Request().Context(), query, prefix, now)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	err = m.Bun.ScanRows(c.Request().Context(), rows, &code)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	return code, nil
}

func (m *psqlReferralsRepository) RSumRefIncentive(c echo.Context, promoCode string, reward models.Reward) (models.SumIncentive, error) {
	var sumIncentive models.SumIncentive
	var perDay, perMonth float64

	query := `select sum(reward_referral) per_day
		from referral_transactions rt
		left join reward_transactions rtrx 
			on rt.used_referral_code = rtrx.used_promo_code 
			and rt.ref_id = rtrx.ref_id
		where rt.used_referral_code = ?0
		and rtrx.transaction_date::date = ?1`

	err := m.Bun.QueryThenScan(c, &perDay, query, promoCode, time.Now().Format(models.DateFormat))

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return sumIncentive, err
	}

	query = `select sum(reward_referral) per_month
		from referral_transactions rt
		left join reward_transactions rtrx 
			on rt.used_referral_code = rtrx.used_promo_code 
			and rt.ref_id = rtrx.ref_id
		where rt.used_referral_code = ?0
		and rtrx.transaction_date between (date_trunc('month', ?1::date)::date) 
			and (((date_trunc('month', ?1::date)) + ('1 month'::INTERVAL))::date)`

	err = m.Bun.QueryThenScan(c, &perMonth, query, promoCode, time.Now())

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return sumIncentive, err
	}

	sumIncentive.Reward = reward
	sumIncentive.PerDay = perDay
	sumIncentive.PerMonth = perMonth

	return sumIncentive, nil
}

func (m *psqlReferralsRepository) RGenerateCode(c echo.Context, refCode models.ReferralCodes, prefix string) string {
	return prefix + helper.RandomStr(5, map[string]bool{})
}

func (m *psqlReferralsRepository) RGetReferralCampaignMetadata(c echo.Context, pv models.PayloadValidator) (models.PrefixResponse, error) {
	var response models.PrefixResponse

	query := `SELECT c.metadata->>'prefix' as prefix FROM campaigns c
		LEFT JOIN rewards r ON c.id = r.campaign_id
		WHERE c.status = ?0 AND c.metadata->>'isReferral' = 'true' AND r.is_promo_code = ?1
		AND c.start_date::date <= ?2 AND (c.end_date::date >= ?2 OR c.end_date IS null)
		order by c.created_at desc limit 1`

	err := m.Bun.QueryThenScan(c, &response.Prefix, query, models.CampaignActive,
		models.IsPromoCodeFalse, pv.TransactionDate)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.PrefixResponse{}, err
	}

	if response.Prefix == "" {
		logger.Make(c, nil).Debug(models.ErrRefPrefixNF)

		return models.PrefixResponse{}, models.ErrRefPrefixNF
	}

	return response, nil
}

func (m *psqlReferralsRepository) RGetHistoryIncentive(c echo.Context, refCif string) ([]models.ResponseHistoryIncentive, error) {
	var historyIncentives []models.ResponseHistoryIncentive

	query := `select
		(rtrx.request_data->>'validators')::json->>'transactionType' as transaction_type, 
		(rtrx.request_data->>'validators')::json->>'product' as product_code,
		rtrx.request_data->>'customerName' as customer_name,
		rt.reward_referral,
		rt.created_at
			from referral_transactions rt 
 			left join reward_transactions rtrx on rtrx.ref_id = rt.ref_id 
 			where rtrx.status = '1'
			and rtrx.request_data->>'referrer' = ?0
 			and rt.reward_type = ?1
 			order by rt.created_at desc;`

	err := m.Bun.QueryThenScan(c, &historyIncentives, query, refCif,
		models.CodeTypeIncentive)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return []models.ResponseHistoryIncentive{}, err
	}

	return historyIncentives, nil
}
