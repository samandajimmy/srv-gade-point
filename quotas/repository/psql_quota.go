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

func (quotRepo *psqlQuotaRepository) Create(c echo.Context, quota *models.Quota, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `INSERT INTO quotas (number_of_days, amount, is_per_user, reward_id, created_at)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`
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

	err = stmt.QueryRow(quota.NumberOfDays, quota.Amount, quota.IsPerUser, rewardID, &now).Scan(&lastID)

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

func (quotRepo *psqlQuotaRepository) CheckQuota(c echo.Context, rewardID int64) (int64, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var available int64

	query := `SELECT available FROM metrics where job = $1`
	stmt, err := quotRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	err = stmt.QueryRow(&rewardID).Scan(&available)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	return available, nil
}
