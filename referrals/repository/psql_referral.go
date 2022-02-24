package repository

import (
	"database/sql"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"time"

	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type psqlReferralsRepository struct {
	Conn *sql.DB
	Bun  *bun.DB
}

// NewPsqlReferralRepository will create an object that represent the referrals.Repository interface
func NewPsqlReferralRepository(Conn *sql.DB, Bun *bun.DB) referrals.Repository {
	return &psqlReferralsRepository{Conn, Bun}
}

func (refRepo *psqlReferralsRepository) PostCoreTrx(c echo.Context, coreTrx models.CoreTrxPayload) error {
	now := time.Now()
	createdAt := now
	trxType := 1

	query := `INSERT INTO referral_transactions
		(cif, ref_id, used_referral_code, type, reward_referral, reward_type, created_at, phone_number,
			trx_amount, loan_amount, interest_amount, trx_id, product_code, trx_date, trx_type)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := refRepo.Bun.QueryContext(c.Request().Context(), query, coreTrx.CIF,
		coreTrx.RefID,
		coreTrx.UsedReferralCode,
		coreTrx.Type,
		coreTrx.RewardReferral,
		coreTrx.RewardType,
		&createdAt,
		coreTrx.PhoneNumber,
		coreTrx.TrxAmount,
		coreTrx.LoanAmount,
		coreTrx.InterestAmount,
		coreTrx.TrxID,
		coreTrx.ProductCode,
		coreTrx.TrxDate,
		trxType)

	if err != nil {
		logger.Make(c, nil).Error(err)
		return err
	}

	return nil
}
