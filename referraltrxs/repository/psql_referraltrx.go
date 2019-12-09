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
