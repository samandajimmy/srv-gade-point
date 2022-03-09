package repository

import (
	"database/sql"
	gcdb "gade/srv-gade-point/database"
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
	stmts := []*gcdb.PipelineStmt{}
	for _, trx := range coreTrx {

		stmts = append(stmts, gcdb.NewPipelineStmt(`INSERT INTO referral_transactions 
		(cif, ref_id, used_referral_code, type, reward_referral, reward_type, created_at, 
		phone_number, trx_amount, loan_amount, interest_amount, product_code, 
		trx_date, trx_type) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)`,
			nilFilters, trx.CIF, trx.RefID, trx.UsedReferralCode, trx.Type, trx.RewardReferral,
			trx.RewardType, createdAt, trx.PhoneNumber, trx.TrxAmount, trx.LoanAmount,
			trx.InterestAmount, trx.ProductCode, trx.TrxDate, trxType))
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
func (m *psqlReferralsRepository) CreateReferral(c echo.Context, refcodes models.ReferralCodes) (models.ReferralCodes, error) {

	now := time.Now()

	query := `INSERT INTO referral_codes (cif, referral_code, campaign_id, created_at) VALUES (?0, ?1, ?2, ?3) RETURNING id`

	_, err := m.Bun.QueryContext(c.Request().Context(), query, refcodes.CIF, refcodes.ReferralCode, refcodes.CampaignId, now)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ReferralCodes{}, err
	}

	refcodes.CreatedAt = now
	refcodes.UpdatedAt = now

	return refcodes, nil
}

func (m *psqlReferralsRepository) GetReferralByCif(c echo.Context, cif string) (models.ReferralCodes, error) {

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

func (m *psqlReferralsRepository) GetCampaignId(c echo.Context, prefix string) (int64, error) {

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

	query := `select sum(reward_referral) per_day
		from referral_transactions rt
		where rt.used_referral_code = ?0 and rt.reward_id = ?1
		and rt.trx_date::date = ?2::date`

	err := m.Bun.QueryThenScan(c, &sumIncentive, query, promoCode, reward.ID, time.Now())

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return sumIncentive, err
	}

	query = `select sum(reward_referral) per_month
		from referral_transactions rt
		where rt.used_referral_code = ?0 and rt.reward_id = ?1
		and trx_date between (date_trunc('month', ?2::date)::date) 
			and (((date_trunc('month', ?2::date)) + ('1 month'::INTERVAL))::date)`

	err = m.Bun.QueryThenScan(c, &sumIncentive, query, promoCode, reward.ID, time.Now())

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return sumIncentive, err
	}

	sumIncentive.Reward = reward

	return sumIncentive, nil
}
