package repository

import (
	"database/sql"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/quotas"
	"time"

	"github.com/labstack/echo"
)

type psqlQuotaRepository struct {
	Conn *sql.DB
}

// NewPsqlQuotaRepository will create an object that represent the quotas.Repository interface
func NewPsqlQuotaRepository(Conn *sql.DB) quotas.Repository {
	return &psqlQuotaRepository{Conn}
}

func (quotRepo *psqlQuotaRepository) Create(c echo.Context, quota *models.Quota, reward *models.Reward) error {
	var nextCheck time.Time
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()

	if *quota.IsPerUser == models.IsPerUserFalse && *quota.NumberOfDays != models.QuotaUnlimited {
		nextCheck, _ = time.Parse(time.RFC3339, reward.Campaign.StartDate)
		nextCheck = nextCheck.AddDate(0, 0, int(*quota.NumberOfDays-1))
	} else {
		nextCheck, _ = time.Parse(time.RFC3339, reward.Campaign.EndDate)
	}

	query := `INSERT INTO quotas (number_of_days, amount, is_per_user, reward_id, available, last_check, next_check, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	stmt, err := quotRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	var lastID int64

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(quota.NumberOfDays, quota.Amount, quota.IsPerUser, reward.ID, quota.Amount, reward.Campaign.StartDate, &nextCheck, &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	quota.ID = lastID
	quota.CreatedAt = &now
	return nil
}

func (quotRepo *psqlQuotaRepository) DeleteByReward(c echo.Context, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `DELETE FROM quotas WHERE reward_id = $1`
	stmt, err := quotRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	result, err := stmt.Query(rewardID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	requestLogger.Debug("quotas deleted: ", result)

	return nil
}

func (quotRepo *psqlQuotaRepository) CheckQuota(c echo.Context, rewardID int64) ([]*models.Quota, error) {
	var result []*models.Quota
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `SELECT id, available, is_per_user, amount, number_of_days, last_check FROM quotas where reward_id = $1`
	rows, err := quotRepo.Conn.Query(query, rewardID)
	defer rows.Close()

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	for rows.Next() {
		var r models.Quota
		var d models.Reward

		err = rows.Scan(
			&r.ID,
			&r.Available,
			&r.IsPerUser,
			&r.Amount,
			&r.NumberOfDays,
			&r.LastCheck,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		d.ID = rewardID
		r.Reward = &d

		result = append(result, &r)
	}

	return result, err
}

func (quotRepo *psqlQuotaRepository) UpdateAddQuota(c echo.Context, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `UPDATE quotas SET available = available + 1 WHERE reward_id = $1 and is_per_user = $2`
	stmt, err := quotRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	_, err = stmt.Query(&rewardID, models.IsPerUserFalse)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func (quotRepo *psqlQuotaRepository) UpdateReduceQuota(c echo.Context, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `UPDATE quotas SET available = available - 1 WHERE reward_id = $1 and is_per_user = $2`
	stmt, err := quotRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	_, err = stmt.Query(&rewardID, models.IsPerUserFalse)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func (quotRepo *psqlQuotaRepository) CheckRefreshQuota(c echo.Context, payload *models.PayloadValidator) ([]*models.Quota, error) {
	var result []*models.Quota
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `select q.id, q.available, q.is_per_user, q.amount, q.number_of_days, q.last_check, q.next_check from quotas q where q.reward_id in (select r.id as reward_id
		from campaigns c inner join rewards r on c.id = r.campaign_id where c.status = $1
		and c.start_date::date < $2 and c.end_date::date > $2)
		and q.is_per_user = $3 and q.number_of_days > $4 and q.amount > $5 and q.next_check::date < $2`
	rows, err := quotRepo.Conn.Query(
		query,
		models.CampaignActive,
		payload.TransactionDate,
		models.IsPerUserFalse,
		models.QuotaUnlimited,
		models.IsLimitAmount,
	)

	defer rows.Close()

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	for rows.Next() {
		var r models.Quota

		err = rows.Scan(
			&r.ID,
			&r.Available,
			&r.IsPerUser,
			&r.Amount,
			&r.NumberOfDays,
			&r.LastCheck,
			&r.NextCheck,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		result = append(result, &r)
	}

	return result, err
}

func (quotRepo *psqlQuotaRepository) RefreshQuota(c echo.Context, qList *models.Quota, payload *models.PayloadValidator) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	transactionDate, _ := time.Parse(models.DateFormat, payload.TransactionDate)
	numberOfDay := int(*qList.NumberOfDays)
	diff := transactionDate.Sub(*qList.LastCheck)
	selisihTanggal := int(diff.Hours() / 24) // number of days
	modulus := selisihTanggal % numberOfDay
	intervalDays := numberOfDay - modulus

	nextCheckDate := transactionDate.AddDate(0, 0, intervalDays-1)
	lastCheckDate := nextCheckDate.AddDate(0, 0, (numberOfDay-1)*-1)

	query := `update quotas set available = $1, last_check = $2, next_check = $3 where id = $4`

	stmt, err := quotRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	_, err = stmt.Query(qList.Amount, lastCheckDate, nextCheckDate, qList.ID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}
