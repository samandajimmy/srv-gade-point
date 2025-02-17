package repository

import (
	"database/sql"
	"fmt"
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

func (m *psqlReferralsRepository) RCreateReferral(c echo.Context, refcodes models.ReferralCodes) (models.ReferralCodes, error) {
	now := time.Now()
	query := `INSERT INTO referral_codes (cif, referral_code, campaign_id, created_at) 
		VALUES (?0, ?1, ?2, ?3) RETURNING id`

	_, err := m.Bun.QueryContext(c.Request().Context(), query, refcodes.CIF, refcodes.ReferralCode, refcodes.CampaignId, now)

	if err != nil {
		logger.Make(c).Debug(err)

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
		logger.Make(c).Debug(err)

		return models.ReferralCodes{}, err
	}

	err = m.Bun.ScanRows(c.Request().Context(), rows, &result)

	if err != nil {
		logger.Make(c).Debug(err)

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
		logger.Make(c).Debug(err)

		return 0, err
	}

	err = m.Bun.ScanRows(c.Request().Context(), rows, &code)

	if err != nil {
		logger.Make(c).Debug(err)

		return 0, err
	}

	return code, nil
}

func (m *psqlReferralsRepository) RSumRefIncentive(c echo.Context, promoCode string) (models.ObjIncentive, error) {
	var objIncentive models.ObjIncentive
	var sums []float64

	query := `select sum(reward_referral) sum
		from referral_transactions rt
		left join reward_transactions rtrx 
			on rt.used_referral_code = rtrx.used_promo_code 
			and rt.ref_id = rtrx.ref_id
		where rt.used_referral_code = ?0
		and rtrx.transaction_date::date = ?1
		UNION ALL
		select sum(reward_referral) sum
		from referral_transactions rt
		left join reward_transactions rtrx 
			on rt.used_referral_code = rtrx.used_promo_code 
			and rt.ref_id = rtrx.ref_id
		where rt.used_referral_code = ?0
		and rtrx.transaction_date between (date_trunc('month', ?1::date)::date) 
			and (((date_trunc('month', ?1::date)) + ('1 month'::INTERVAL))::date)
		`

	err := m.Bun.QueryThenScan(c, &sums, query, promoCode, time.Now().Format(models.DateFormat))

	if err != nil {
		logger.Make(c).Debug(err)

		return objIncentive, err
	}

	objIncentive.PerDay = sums[0]
	objIncentive.PerMonth = sums[1]

	return objIncentive, nil
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
		logger.Make(c).Debug(err)

		return models.PrefixResponse{}, err
	}

	if response.Prefix == "" {
		logger.Make(c).Debug(models.ErrRefPrefixNF)

		return models.PrefixResponse{}, models.ErrRefPrefixNF
	}

	return response, nil
}

func (m *psqlReferralsRepository) RGetHistoryIncentive(c echo.Context, pl models.RequestHistoryIncentive) (models.ResponseHistoryIncentive, error) {
	var historyIncentives models.ResponseHistoryIncentive
	var historyData []models.ResponseHistoryIncentiveData
	var totalData int64

	query := `select
		(rtrx.request_data->>'validators')::json->>'transactionType' as transaction_type, 
		(rtrx.request_data->>'validators')::json->>'product' as product_code,
		rtrx.request_data->>'customerName' as customer_name,
		rt.reward_referral,
		EXTRACT(EPOCH FROM rt.created_at::timestamp) as created_at
			from referral_transactions rt 
 			left join reward_transactions rtrx on rtrx.ref_id = rt.ref_id 
 			where rtrx.status = '1'
			and rtrx.request_data->>'referrer' = ?0
 			and rt.reward_type = ?1
 			order by rt.created_at desc limit ?2`

	if pl.Page > 0 {
		paging := fmt.Sprintf(" OFFSET %d", ((pl.Page - 1) * pl.Limit))
		query += paging
	}

	err := m.Bun.QueryThenScan(c, &historyData, query, pl.RefCif,
		models.CodeTypeIncentive, pl.Limit)

	if err != nil {
		logger.Make(c).Debug(err)

		return models.ResponseHistoryIncentive{}, err
	}

	query = `select count(rt.id) as total_data
			from referral_transactions rt 
 			left join reward_transactions rtrx on rtrx.ref_id = rt.ref_id 
 				where rtrx.request_data->>'referrer' = ? 
 				and rtrx.status = '1'
 				and rt.reward_type = 'incentive';`

	err = m.Bun.QueryThenScan(c, &totalData, query, pl.RefCif)

	if err != nil {
		logger.Make(c).Debug(err)

		return models.ResponseHistoryIncentive{}, err
	}

	historyIncentives.TotalData = totalData
	historyIncentives.HistoryIncentiveData = &historyData

	return historyIncentives, nil
}

func (m *psqlReferralsRepository) RTotalFriends(c echo.Context, cif string) (models.RespTotalFriends, error) {
	var friends []models.Friends
	var resTFriend models.RespTotalFriends

	query := `SELECT
			DISTINCT customer_name
			FROM (
				SELECT rtrx.request_data ->>'customerName' as customer_name
				from referral_transactions rt 
				left join reward_transactions rtrx on rtrx.ref_id = rt.ref_id 
				where rtrx.request_data->>'referrer' = ? 
				and rtrx.status = '1'
				and rt.reward_type = 'incentive'
			) AS subquery;`

	err := m.Bun.QueryThenScan(c, &friends, query, cif)

	if err != nil {
		logger.Make(c).Debug(err)

		return models.RespTotalFriends{}, err
	}

	resTFriend.TotalFriends = len(friends)

	return resTFriend, nil
}

func (m *psqlReferralsRepository) RFriendsReferral(c echo.Context, pl models.PayloadFriends) ([]models.Friends, error) {

	var refMembers []models.Friends

	query := `select customer_name from (SELECT DISTINCT(rtrx.request_data ->>'customerName') as customer_name,
	rtrx.transaction_date as transaction_date
	from referral_transactions rt 
	left join reward_transactions rtrx on rtrx.ref_id = rt.ref_id 
	where rtrx.request_data->>'referrer' = ? 
	and rtrx.status = '1'
	and rt.reward_type = 'incentive'
	LIMIT ?`

	paging := ""

	if pl.Page > 0 {
		paging = fmt.Sprintf(" OFFSET %d", ((pl.Page - 1) * pl.Limit))
	}

	orderQuery := ") s order by transaction_date desc"
	query += paging + orderQuery

	err := m.Bun.QueryThenScan(c, &refMembers, query, pl.CIF, pl.Limit)

	if err != nil {
		logger.Make(c).Debug(err)
		return nil, err
	}

	return refMembers, nil
}

func (m *psqlReferralsRepository) RGetReferralCodeByCampaignId(c echo.Context, campaignId int64) ([]string, error) {
	var arrRef []string

	query := `select referral_code
	from referral_codes
 	where campaign_id = ?0 ;`

	err := m.Bun.QueryThenScan(c, &arrRef, query, campaignId)

	if err != nil {
		logger.Make(c).Debug(err)

		return arrRef, err
	}

	return arrRef, nil
}
