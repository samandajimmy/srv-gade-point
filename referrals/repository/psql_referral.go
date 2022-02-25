package repository

import (
	"database/sql"
	"fmt"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"strings"
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

func (refRepo *psqlReferralsRepository) PostCoreTrx(c echo.Context, coreTrx []models.CoreTrxPayload) error {
	const fInsertTrx = 15
	now := time.Now()
	createdAt := now
	trxType := 1
	i := 0
	valueArgs := []interface{}{}
	valueStrings := []string{}

	for _, trx := range coreTrx {
		genIndex := "("
		genStartNumber := i * fInsertTrx
		genLastNumber := genStartNumber + fInsertTrx

		for x := genStartNumber; x < genLastNumber; x++ {
			genIndex += fmt.Sprintf("$%d", x+1)
			if x != (genLastNumber - 1) {
				genIndex += ", "
			}
		}

		genIndex += ")"
		valueStrings = append(valueStrings, genIndex)

		valueArgs = append(valueArgs, trx.CIF)
		valueArgs = append(valueArgs, trx.RefID)
		valueArgs = append(valueArgs, trx.UsedReferralCode)
		valueArgs = append(valueArgs, trx.Type)
		valueArgs = append(valueArgs, trx.RewardReferral)
		valueArgs = append(valueArgs, trx.RewardType)
		valueArgs = append(valueArgs, createdAt)
		valueArgs = append(valueArgs, trx.PhoneNumber)
		valueArgs = append(valueArgs, trx.TrxAmount)
		valueArgs = append(valueArgs, trx.LoanAmount)
		valueArgs = append(valueArgs, trx.InterestAmount)
		valueArgs = append(valueArgs, trx.TrxID)
		valueArgs = append(valueArgs, trx.ProductCode)
		valueArgs = append(valueArgs, trx.TrxDate)
		valueArgs = append(valueArgs, trxType)
		i++
	}

	query := fmt.Sprintf(`INSERT INTO referral_transactions 
		(cif, ref_id, used_referral_code, type, reward_referral, reward_type, created_at, phone_number,
			trx_amount, loan_amount, interest_amount, trx_id, product_code, trx_date, trx_type) VALUES %s`, strings.Join(valueStrings, ","))
	stmt, err := refRepo.Conn.Prepare(query)

	if err != nil {
		logger.Make(c, nil).Error(err)
		return err
	}

	rows, err := stmt.Query(valueArgs...)

	if err != nil {
		logger.Make(c, nil).Error(err)
		return err
	}

	defer rows.Close()
	logger.Make(c, nil).Debug("voucher code(s) are created concurrently!")

	return nil
}
