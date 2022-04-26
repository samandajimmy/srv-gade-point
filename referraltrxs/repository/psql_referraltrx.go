package repository

import (
	"database/sql"
	"fmt"
	gcdb "gade/srv-gade-point/database"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referraltrxs"
	"time"

	"github.com/labstack/echo"
)

type psqlReferralTrxRepository struct {
	Conn *sql.DB
	Bun  *gcdb.DbBun
}

// NewPsqlReferralTrxRepository will create an object that represent the referraltrxs.Repository interface
func NewPsqlReferralTrxRepository(Conn *sql.DB, Bun *gcdb.DbBun) referraltrxs.RefTRepository {
	return &psqlReferralTrxRepository{Conn, Bun}
}

func (refTrxRepo *psqlReferralTrxRepository) RPostReferralTrx(c echo.Context, refTrx models.ReferralTrx) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `INSERT INTO referral_transactions (cif, ref_id, used_referral_code, type,
		reward_referral, reward_type, phone_number, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`

	lastID := int64(0)
	stmt, err := refTrxRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(
		&refTrx.CIF, &refTrx.RefID, &refTrx.UsedReferralCode, &refTrx.Type, &refTrx.RewardReferral,
		&refTrx.RewardType, &refTrx.PhoneNumber, &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func (refTrxRepo *psqlReferralTrxRepository) RGetMilestone(c echo.Context, payload models.MilestonePayload) (*models.Milestone, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	result := new(models.Milestone)

	query := fmt.Sprintf(`SELECT count(r.id) as totalRewardCounter, sum(r.reward_referral) AS totalReward
			  FROM referral_transactions r
			  WHERE used_referral_code = '%s' and type = '%d'`,
		payload.ReferralCode, models.ReferralType[models.RefTargetReferrer])

	err := refTrxRepo.Conn.QueryRow(query).Scan(
		&result.TotalRewardCounter,
		&result.TotalReward,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return result, nil
}

func (refTrxRepo *psqlReferralTrxRepository) RGetRanking(c echo.Context, rp models.RankingPayload) ([]*models.Ranking, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var prefixReferral = "PDS%"
	query := fmt.Sprintf(`(select used, used_referral_code, date, ROW_NUMBER () OVER (ORDER BY used desc) as row,
		CASE WHEN used_referral_code = '%s' THEN true
		ELSE false
		END as is_selected
		from (
			select count(used_referral_code) as used, used_referral_code, max(created_at) as date
			from referral_transactions
			where type = 1 and created_at between '%s' and '%s' and used_referral_code LIKE '%s'
			group by used_referral_code
			order by used desc, date asc
		) as bro limit 10)
		union
		select used, used_referral_code, date, row, true as is_selected
		from (
			select used_referral_code, used, date, ROW_NUMBER () OVER (ORDER BY used desc) as row from (
				select count(used_referral_code) as used, used_referral_code, max(created_at) as date
				from referral_transactions
				where type = 1 and created_at between '%s' and '%s' and used_referral_code LIKE '%s'
				group by used_referral_code
				order by used desc, date asc
			) foo
			group by used_referral_code, used, date
			order by used desc, date asc
		) as bro
		where used_referral_code = '%s'
		order by used desc, date asc`, rp.ReferralCode, rp.StartDate, rp.EndDate, prefixReferral, rp.StartDate, rp.EndDate, prefixReferral, rp.ReferralCode)

	rows, err := refTrxRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	result := make([]*models.Ranking, 0)

	for rows.Next() {
		ranking := new(models.Ranking)
		err = rows.Scan(
			&ranking.TotalUsed,
			&ranking.ReferralCode,
			&ranking.Date,
			&ranking.NoRanking,
			&ranking.IsReferralCode,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		result = append(result, ranking)
	}

	return result, nil
}

func (refTrxRepo *psqlReferralTrxRepository) RGetRankingByReferralCode(c echo.Context, referralCode string) (*models.Ranking, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	result := new(models.Ranking)
	query := fmt.Sprintf(`select topTen.* 
		from (select used_referral_code,
		count(used_referral_code) as total,
		row_number() over (order by count(used_referral_code) desc) as URUT
		from referral_transactions
		where type = '1' AND created_at >= date_trunc('month', CURRENT_DATE)
		group by used_referral_code
		order by total desc
		limit 10
		offset 0) as topTen
		where topTen.used_referral_code = '%s'`, referralCode)

	err := refTrxRepo.Conn.QueryRow(query, referralCode).Scan(
		&result.ReferralCode,
		&result.TotalUsed,
		&result.NoRanking,
	)

	if err != nil && err != sql.ErrNoRows {
		requestLogger.Debug(err)

		return nil, err
	}

	return result, nil
}

func (refTrxRepo *psqlReferralTrxRepository) RGetCountRefTrxByCIF(c echo.Context, cif string, productCode string, promoCode string) bool {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var totalPromo int64

	query := `SELECT COUNT(ref.used_referral_code) 
	FROM referral_transactions ref	
	LEFT JOIN reward_transactions rwd ON ref.CIF = rwd.CIF AND rwd.used_promo_code = ?0
	WHERE ref.CIF = ?1
	AND rwd.succeeded_date IS NOT NULL
	AND rwd.request_data->'validators'->>'product' = ?2`

	err := refTrxRepo.Bun.QueryThenScan(c, &totalPromo, query, promoCode, cif, productCode)

	if err != nil {
		requestLogger.Debug(err)

		return false
	}

	if totalPromo >= 1 {
		return false
	}

	return true
}
