package repository

import (
	"database/sql"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referraltrxs"
	"time"

	"github.com/labstack/echo"
)

type psqlReferralTrxRepository struct {
	Conn *sql.DB
}

// NewPsqlReferralTrxRepository will create an object that represent the referraltrxs.Repository interface
func NewPsqlReferralTrxRepository(Conn *sql.DB) referraltrxs.Repository {
	return &psqlReferralTrxRepository{Conn}
}

func (refTrxRepo *psqlReferralTrxRepository) CheckReferralTrxByExistingReferrer(c echo.Context, refTrx models.ReferralTrx) (int64, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT COUNT(ref.id) as Count FROM referral_transactions ref 
			where ref.cif = $1 and ref.type = 0 `
	count := int64(0)
	err := refTrxRepo.Conn.QueryRow(query, refTrx.CifReferrer).Scan(&count)

	if err != nil {
		requestLogger.Debug(err)
		return 0, err
	}

	return count, nil
}

func (refTrxRepo *psqlReferralTrxRepository) CheckReferralTrxByValueRewards(c echo.Context, refTrx models.ReferralTrx) (*models.ReferralTrx, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	result := new(models.ReferralTrx)
	query := `SELECT COALESCE(SUM(ref.reward_referral), 0) FROM referral_transactions ref 
			where ref.cif = $1 and type = 1`
	err := refTrxRepo.Conn.QueryRow(query, refTrx.CIF).Scan(
		&result.TotalGoldback)

	if err != nil {
		requestLogger.Debug(err)
		return nil, err
	}

	return result, nil
}

func (refTrxRepo *psqlReferralTrxRepository) Create(c echo.Context, refTrx models.ReferralTrx) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `INSERT INTO referral_transactions (cif, ref_id, used_referral_code, type,
		reward_referral, reward_type, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`

	lastID := int64(0)
	stmt, err := refTrxRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(
		&refTrx.CIF, &refTrx.RefID, &refTrx.UsedReferralCode, &refTrx.Type, &refTrx.RewardReferral,
		&refTrx.RewardType, &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func (refTrxRepo *psqlReferralTrxRepository) GetMilestone(c echo.Context, CIF *int64) (*models.Milestone, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	result := new(models.Milestone)

	query := `SELECT count(r.id) as totalRewardCounter, sum(r.reward_referral) AS totalReward
			  FROM referral_transactions r
			  WHERE cif = $1 and type = 1`

	err := refTrxRepo.Conn.QueryRow(query, CIF).Scan(
		&result.TotalRewardCounter,
		&result.TotalReward,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return result, nil
}

func (refTrxRepo *psqlReferralTrxRepository) GetRanking(c echo.Context, referralCode string) ([]models.Ranking, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	isInRanking := false
	query := `select topTen.*
		from (select used_referral_code as referral_code,
		count(used_referral_code) as total,
		row_number() over (order by count(used_referral_code) desc) as rank
		from referral_transactions
		where type = '1' AND created_at >= date_trunc('month', CURRENT_DATE)
		group by used_referral_code
		order by total desc
		limit 10
		offset 0) as topTen`

	rows, err := refTrxRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()
	var result []models.Ranking
	for rows.Next() {
		var ranking models.Ranking
		err = rows.Scan(
			&ranking.ReferralCode,
			&ranking.TotalUsed,
			&ranking.NoRanking,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		if referralCode == ranking.ReferralCode {
			isInRanking = true
		}

		result = append(result, ranking)
	}

	if !isInRanking {
		ranking, err := refTrxRepo.GetRankingByReferralCode(c, referralCode)
		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}
		result = append(result, *ranking)
	}

	return result, nil
}

func (refTrxRepo *psqlReferralTrxRepository) GetRankingByReferralCode(c echo.Context, referralCode string) (*models.Ranking, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	result := new(models.Ranking)
	query := `select topTen.* 
		from (select used_referral_code,
		count(used_referral_code) as total,
		row_number() over (order by count(used_referral_code) desc) as URUT
		from referral_transactions
		where type = '1' AND created_at >= date_trunc('month', CURRENT_DATE)
		group by used_referral_code
		order by total desc
		limit 10
		offset 0) as topTen
		where topTen.used_referral_code = $1`

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
