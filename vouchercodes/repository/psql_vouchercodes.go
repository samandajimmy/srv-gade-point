package repository

import (
	"database/sql"
	"fmt"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchercodes"

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
