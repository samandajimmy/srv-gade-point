package repository

import (
	"database/sql"
	"fmt"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchercodes"
	"time"

	"github.com/labstack/echo"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type psqlVoucherCodeRepository struct {
	Conn *sql.DB
}

// NewPsqlVoucherCodeRepository will create an object that represent the vouchercode.Repository interface
func NewPsqlVoucherCodeRepository(Conn *sql.DB) vouchercodes.Repository {
	return &psqlVoucherCodeRepository{Conn}
}

func (psqlRepo *psqlVoucherCodeRepository) CountVoucherCode(c echo.Context, payload map[string]interface{}) (string, error) {
	var counter string
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	userID := payload["userId"].(string)
	where := ""

	queryCounter := `SELECT COUNT(ID) counter FROM voucher_codes where user_id is not null`

	if userID != "" {
		where += " and user_id = '" + userID + "'"
	}

	queryCounter += where + ";"
	err := psqlRepo.Conn.QueryRow(queryCounter).Scan(&counter)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}

	return counter, nil
}

func (psqlRepo *psqlVoucherCodeRepository) GetVoucherCodeHistory(c echo.Context, payload map[string]interface{}) ([]models.VoucherCode, error) {
	var result []models.VoucherCode
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	user := payload["userId"].(string)
	paging := ""
	where := ""

	query := `SELECT vc.id, coalesce(vc.user_id, ''), vc.promo_code, vc.status, vc.bought_date, vc.redeemed_date, vc.updated_at, v.id, v.name
			FROM voucher_codes vc left join vouchers v on vc.voucher_id = v.id where user_id is not null`

	if user != "" {
		where += " and user_id = '" + user + "'"
	}

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	query += where + " order by vc.updated_at desc, status desc" + paging + ";"
	rows, err := psqlRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		vchrCode := models.VoucherCode{}
		voucher := models.Voucher{}

		err = rows.Scan(
			&vchrCode.ID,
			&vchrCode.UserID,
			&vchrCode.PromoCode,
			&vchrCode.Status,
			&vchrCode.BoughtDate,
			&vchrCode.RedeemedDate,
			&vchrCode.UpdatedAt,
			&voucher.ID,
			&voucher.Name,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		if vchrCode.ID != 0 {
			vchrCode.Voucher = &voucher
		}

		result = append(result, vchrCode)
	}

	return result, nil
}

func (psqlRepo *psqlVoucherCodeRepository) CountVoucherCodeByVoucherID(c echo.Context, payload map[string]interface{}) (string, error) {
	var counter string
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	voucherID := payload["voucherId"].(string)
	where := ""

	queryCounter := `SELECT COUNT(ID) counter FROM voucher_codes`

	if voucherID != "" {
		where += " where voucher_id = '" + voucherID + "'"
	}

	queryCounter += where + ";"
	err := psqlRepo.Conn.QueryRow(queryCounter).Scan(&counter)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}

	return counter, nil
}

func (psqlRepo *psqlVoucherCodeRepository) GetVoucherCodes(c echo.Context, payload map[string]interface{}) ([]models.VoucherCode, error) {
	var result []models.VoucherCode
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	voucherID := payload["voucherId"].(string)
	paging := ""
	where := ""

	query := `SELECT vc.id, vc.promo_code, vc.status, vc.bought_date, vc.redeemed_date, v.id, v.name, v.end_date
	        FROM voucher_codes vc left join vouchers v on vc.voucher_id = v.id`

	if voucherID != "" {
		where += " where voucher_id = '" + voucherID + "'"
	}

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	query += where + " order by vc.id asc" + paging + ";"
	rows, err := psqlRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		vchrCode := models.VoucherCode{}
		voucher := models.Voucher{}

		err = rows.Scan(
			&vchrCode.ID,
			&vchrCode.PromoCode,
			&vchrCode.Status,
			&vchrCode.BoughtDate,
			&vchrCode.RedeemedDate,
			&voucher.ID,
			&voucher.Name,
			&voucher.EndDate,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		if vchrCode.ID != 0 {
			vchrCode.Voucher = &voucher
		}

		result = append(result, vchrCode)
	}

	return result, nil
}

func (psqlRepo *psqlVoucherCodeRepository) GetVoucherCodeRefID(c echo.Context, refID string) (*models.VoucherCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `SELECT vc.promo_code, v.name FROM voucher_codes vc
		left join vouchers v on vc.voucher_id = v.id where ref_id = $1`

	stmt, err := psqlRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	vchrCode := models.VoucherCode{}
	voucher := models.Voucher{}
	err = stmt.QueryRow(&refID).Scan(
		&vchrCode.PromoCode,
		&voucher.Name,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	vchrCode.Voucher = &voucher

	return &vchrCode, nil
}

func (psqlRepo *psqlVoucherCodeRepository) UpdateVoucherCodeRedeemed(c echo.Context, redeemDate string, userID string, promoCode string) (*models.VoucherCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	result := new(models.VoucherCode)
	where := ""

	queryUpdate := `UPDATE voucher_codes SET status = 2, redeemed_date = $1, updated_at = $2 `

	if userID != "" && promoCode != "" {
		where += "where user_id = '" + userID + "' and promo_code = '" + promoCode + "' AND"
	}

	queryUpdate += where + " status = 1 RETURNING promo_code, redeemed_date;"
	stmt, err := psqlRepo.Conn.Prepare(queryUpdate)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	err = stmt.QueryRow(&redeemDate, &now).Scan(&result.PromoCode, &result.RedeemedDate)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return result, nil
}

func (psqlRepo *psqlVoucherCodeRepository) CountBoughtVoucherCode(c echo.Context, payload map[string]interface{}) (string, error) {
	var counter string
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	userID := payload["userId"].(string)
	promoCode := payload["promoCode"].(string)
	where := "where"

	queryCounter := `SELECT COUNT(ID) counter FROM voucher_codes `

	if userID != "" && promoCode == "" {
		where += " user_id = '" + userID + "' AND"
	}

	if promoCode != "" && userID == "" {
		where += " promo_code = '" + promoCode + "' AND"
	}

	if userID != "" && promoCode != "" {
		where += " user_id = '" + userID + "' and promo_code = '" + promoCode + "' AND"
	}

	queryCounter += where + " user_id is not null;"
	err := psqlRepo.Conn.QueryRow(queryCounter).Scan(&counter)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}

	return counter, nil
}

func (psqlRepo *psqlVoucherCodeRepository) GetBoughtVoucherCode(c echo.Context, payload map[string]interface{}) ([]models.VoucherCode, error) {
	var result []models.VoucherCode
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	userID := payload["userId"].(string)
	promoCode := payload["promoCode"].(string)
	where := "where"
	paging := ""

	query := `SELECT id, coalesce(user_id, ''), promo_code, status, bought_date, redeemed_date FROM voucher_codes `

	if userID != "" && promoCode == "" {
		where += " user_id = '" + userID + "' AND"
	}

	if promoCode != "" && userID == "" {
		where += " promo_code = '" + promoCode + "' AND"
	}

	if userID != "" && promoCode != "" {
		where += " user_id = '" + userID + "' and promo_code = '" + promoCode + "' AND"
	}

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	query += where + " user_id is not null order by updated_at desc, status desc " + paging + ";"
	rows, err := psqlRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		vchrCode := models.VoucherCode{}

		err = rows.Scan(
			&vchrCode.ID,
			&vchrCode.UserID,
			&vchrCode.PromoCode,
			&vchrCode.Status,
			&vchrCode.BoughtDate,
			&vchrCode.RedeemedDate,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		result = append(result, vchrCode)

	}

	return result, nil
}

func (psqlRepo *psqlVoucherCodeRepository) UpdateVoucherCodeRefID(c echo.Context, voucherCode *models.VoucherCode, refID string) error {
	if voucherCode == nil {
		return nil
	}

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `UPDATE voucher_codes SET ref_id = $1, updated_at = $2 where status = 1 and id = $3`
	stmt, err := psqlRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	rows, err := stmt.Query(&refID, &now, &voucherCode.ID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	defer rows.Close()

	return nil
}

func (psqlRepo *psqlVoucherCodeRepository) UpdateVoucherCodeRejected(c echo.Context, refID string) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `UPDATE voucher_codes SET status = $1, user_id = $2, bought_date = NULL, ref_id = ''
		updated_at = $3 where status = $4 and ref_id = $5`
	stmt, err := psqlRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	rows, err := stmt.Query(&models.VoucherCodeStatusAvailable, "", &now,
		&models.VoucherCodeStatusBooked, &refID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	defer rows.Close()

	return nil
}

func (psqlRepo *psqlVoucherCodeRepository) UpdateVoucherCodeSucceeded(c echo.Context, refID string) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `UPDATE voucher_codes SET status = $1, bought_date = $2, updated_at = $3
		where status = $4 and ref_id = $5`
	stmt, err := psqlRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	rows, err := stmt.Query(&models.VoucherCodeStatusBought, &now, &now,
		models.VoucherCodeStatusBooked, &refID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	defer rows.Close()

	return nil
}

func (psqlRepo *psqlVoucherCodeRepository) ValidateVoucherGive(c echo.Context, payloadVoucherBuy *models.PayloadVoucherBuy) (*models.VoucherCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `SELECT vc.promo_code, v.name, vc.ref_id FROM voucher_codes vc left join vouchers
		v on vc.voucher_id = v.id where ref_id = $1`

	stmt, err := psqlRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	vchrCode := models.VoucherCode{}
	voucher := models.Voucher{}
	err = stmt.QueryRow(payloadVoucherBuy.RefID).Scan(
		&vchrCode.PromoCode,
		&voucher.Name,
		&vchrCode.RefID,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	vchrCode.Voucher = &voucher

	return &vchrCode, nil
}
