package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat           = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	fieldInsertPromoCode = 3
)

var (
	count = []int{1, 2, 3}
)

type psqlVoucherRepository struct {
	Conn *sql.DB
}

// NewPsqlVoucherRepository will create an object that represent the vouchers. Repository interface
func NewPsqlVoucherRepository(Conn *sql.DB) vouchers.Repository {
	return &psqlVoucherRepository{Conn}
}

func (m *psqlVoucherRepository) CreateVoucher(c echo.Context, voucher *models.Voucher) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	var lastID int64
	query := `INSERT INTO vouchers (name, description, start_date, end_date, point, journal_account, value, 
		image_url, status, stock, prefix_promo_code, validators, terms_and_conditions, how_to_use, limit_per_user, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)  RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	validator, err := json.Marshal(voucher.Validators)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(voucher.Name, voucher.Description, voucher.StartDate, voucher.EndDate,
		voucher.Point, voucher.JournalAccount, voucher.Value, voucher.ImageURL, voucher.Status, voucher.Stock, voucher.PrefixPromoCode,
		string(validator), voucher.TermsAndConditions, voucher.HowToUse, voucher.LimitPerUser, &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	voucher.ID = lastID
	voucher.CreatedAt = &now
	return nil
}

func (m *psqlVoucherRepository) CreatePromoCode(c echo.Context, promoCodes []*models.PromoCode) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var valueStrings []string
	var valueArgs []interface{}
	i := 0

	for _, promoCode := range promoCodes {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d)", i*fieldInsertPromoCode+count[0], i*fieldInsertPromoCode+count[1], i*fieldInsertPromoCode+count[2]))
		valueArgs = append(valueArgs, promoCode.PromoCode)
		valueArgs = append(valueArgs, promoCode.Voucher.ID)
		valueArgs = append(valueArgs, promoCode.CreatedAt)
		i++
	}

	query := fmt.Sprintf("INSERT INTO promo_codes (promo_code, voucher_id, created_at) VALUES %s", strings.Join(valueStrings, ","))
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	result, err := stmt.Query(valueArgs...)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	requestLogger.Debug("Result created promo code: ", result)
	return nil
}

func (m *psqlVoucherRepository) UpdateVoucher(c echo.Context, id int64, updateVoucher *models.UpdateVoucher) error {
	now := time.Now()
	var lastID int64
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `UPDATE vouchers SET status = $1, updated_at = $2 WHERE id = $3 RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(updateVoucher.Status, &now, id).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func (m *psqlVoucherRepository) GetVouchersAdmin(c echo.Context, payload map[string]interface{}) ([]*models.Voucher, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	paging := ""
	where := ""
	query := `SELECT d.id, d.name, d.description, d.start_date, d.end_date, d.point, d.journal_account, d.value, d.image_url, d.status, d.stock, d.prefix_promo_code, d.amount, 
	CASE WHEN d.end_date::date < now()::date THEN 0 ELSE coalesce(e.available, 0) END AS available, coalesce(f.bought, 0) bought, coalesce(g.reedem, 0) reedem, 
	CASE WHEN coalesce(h.expired, 0) - coalesce(g.reedem, 0) < 0 THEN 0 ELSE coalesce(h.expired, 0) - coalesce(g.reedem, 0) END AS expired, d.validators, 
	d.terms_and_conditions, d.how_to_use, d.limit_per_user, d.updated_at, d.created_at 
	FROM (SELECT b.id, b.name, b.description, b.start_date, b.end_date, b.point, b.journal_account, b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, 
	count(a.id) as amount, b.validators, b.terms_and_conditions, b.how_to_use, b.limit_per_user, b.updated_at, b.created_at 
	FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id GROUP BY b.id, b.name, b.description, b.start_date, b.end_date, b.point, b.journal_account, 
	b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, b.validators, b.terms_and_conditions, b.how_to_use, b.limit_per_user, b.updated_at, b.created_at) as d 
	LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as available 
	FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id WHERE a.status = 0 GROUP BY b.id) as e ON e.id = d.id LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as bought 
	FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id WHERE a.status = 1 GROUP BY b.id) as f ON f.id = d.id LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as reedem 
	FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id WHERE a.status = 2 GROUP BY b.id) as g ON g.id = d.id LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as expired 
	FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id WHERE end_date::date < now()::date GROUP BY b.id) as h ON h.id = d.id WHERE d.id IS NOT NULL`

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	if payload["name"].(string) != "" {
		where += " AND d.name LIKE '%" + payload["name"].(string) + "%'"
	}

	if payload["status"].(string) != "" {
		where += " AND d.status='" + payload["status"].(string) + "'"
	}

	if payload["startDate"].(string) != "" {
		where += " AND d.start_date::timestamp::date >= '" + payload["startDate"].(string) + "'"
	}

	if payload["endDate"].(string) != "" {
		where += " AND d.end_date::timestamp::date <= '" + payload["endDate"].(string) + "'"
	}

	query += where + " ORDER BY d.created_at DESC " + paging
	res, err := m.getVouchersAdmin(c, query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return res, nil

}

func (m *psqlVoucherRepository) getVouchersAdmin(c echo.Context, query string) ([]*models.Voucher, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var validator json.RawMessage
	rows, err := m.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()
	result := make([]*models.Voucher, 0)

	for rows.Next() {
		t := new(models.Voucher)
		var createDate, updateDate pq.NullTime
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.StartDate,
			&t.EndDate,
			&t.Point,
			&t.JournalAccount,
			&t.Value,
			&t.ImageURL,
			&t.Status,
			&t.Stock,
			&t.PrefixPromoCode,
			&t.Amount,
			&t.Available,
			&t.Bought,
			&t.Redeemed,
			&t.Expired,
			&validator,
			&t.TermsAndConditions,
			&t.HowToUse,
			&t.LimitPerUser,
			&updateDate,
			&createDate,
		)

		t.CreatedAt = &createDate.Time
		t.UpdatedAt = &updateDate.Time
		err = json.Unmarshal([]byte(validator), &t.Validators)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *psqlVoucherRepository) GetVoucherAdmin(c echo.Context, voucherID string) (*models.Voucher, error) {
	var validator json.RawMessage
	var createDate, updateDate pq.NullTime
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	result := new(models.Voucher)
	query := `SELECT d.id, d.name, d.description, d.start_date, d.end_date, d.point, d.journal_account, d.value, d.image_url, d.status, d.stock, d.prefix_promo_code, d.amount, 
	CASE WHEN d.end_date::date < now()::date THEN 0 ELSE coalesce(e.available, 0) END AS available, coalesce(f.bought, 0) bought , coalesce(g.reedem, 0) reedem, 
	CASE WHEN coalesce(h.expired, 0) - coalesce(g.reedem, 0) < 0 THEN 0 ELSE coalesce(h.expired, 0) - coalesce(g.reedem, 0) END AS expired, d.validators, d.terms_and_conditions, d.how_to_use, 
	d.limit_per_user, d.updated_at, d.created_at FROM (SELECT b.id, b.name, b.description, b.start_date, b.end_date, b.point, b.journal_account, b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, 
	count(a.id) as amount, b.validators, b.terms_and_conditions, b.how_to_use, b.limit_per_user, b.updated_at, b.created_at FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id GROUP BY b.id, b.name, 
	b.description, b.start_date, b.end_date, b.point, b.journal_account, b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, b.validators, b.updated_at, b.created_at) as d 
	LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as available FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id WHERE a.status = 0 GROUP BY b.id) as e ON e.id = d.id 
	LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as bought FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id WHERE a.status = 1 GROUP BY b.id) as f ON f.id = d.id 
	LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as reedem FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id WHERE a.status = 2 GROUP BY b.id) as g ON g.id = d.id 
	LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as expired FROM promo_codes a LEFT JOIN vouchers b ON b.id = a.voucher_id WHERE end_date::date < now()::date GROUP BY b.id) as h ON h.id = d.id 
	WHERE d.id = $1`

	err := m.Conn.QueryRow(query, voucherID).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.StartDate,
		&result.EndDate,
		&result.Point,
		&result.JournalAccount,
		&result.Value,
		&result.ImageURL,
		&result.Status,
		&result.Stock,
		&result.PrefixPromoCode,
		&result.Amount,
		&result.Available,
		&result.Bought,
		&result.Redeemed,
		&result.Expired,
		&validator,
		&result.TermsAndConditions,
		&result.HowToUse,
		&result.LimitPerUser,
		&updateDate,
		&createDate)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	result.CreatedAt = &createDate.Time
	result.UpdatedAt = &updateDate.Time
	err = json.Unmarshal([]byte(validator), &result.Validators)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return result, err
}

func (m *psqlVoucherRepository) GetVouchers(c echo.Context, payload map[string]interface{}) ([]*models.Voucher, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	paging := ""
	where := ""
	query := `SELECT c.id, c.name, c.description, c.start_date, c.end_date, c.point, c.value, c.image_url, c.stock, coalesce(d.available, 0), c.terms_and_conditions, c.how_to_use, c.limit_per_user
	FROM vouchers c LEFT JOIN(SELECT b.id, coalesce(count(a.id), 0) as available FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 0 GROUP BY b.id) d 
	ON d.id = c.id WHERE c.status = 1 AND c.end_date::date >= now()`

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	if payload["name"].(string) != "" {
		where += " AND d.name LIKE '%" + payload["name"].(string) + "%'"
	}

	if payload["startDate"].(string) != "" {
		where += " AND d.start_date::timestamp::date >= '" + payload["startDate"].(string) + "'"
	}

	if payload["endDate"].(string) != "" {
		where += " AND d.end_date::timestamp::date <= '" + payload["endDate"].(string) + "'"
	}

	query += where + " ORDER BY c.created_at DESC " + paging
	rows, err := m.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()
	result := make([]*models.Voucher, 0)

	for rows.Next() {
		t := new(models.Voucher)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.StartDate,
			&t.EndDate,
			&t.Point,
			&t.Value,
			&t.ImageURL,
			&t.Stock,
			&t.Available,
			&t.TermsAndConditions,
			&t.HowToUse,
			&t.LimitPerUser,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *psqlVoucherRepository) UpdateExpiryDate(c echo.Context) error {
	now := time.Now()
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `UPDATE vouchers SET status = 0, updated_at = $1 WHERE end_date::timestamp::date < now()::date AND status = 1`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug("Update Status Base on Expiry Date: ", err)

		return err
	}

	var lastID int64
	err = stmt.QueryRow(&now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug("Update Status Base on Expiry Date: ", err)

		return err
	}

	return nil
}

func (m *psqlVoucherRepository) UpdateStatusBasedOnStartDate() error {
	now := time.Now()
	query := `UPDATE vouchers SET status = 1, updated_at = $1 WHERE start_date::timestamp::date = now()::date`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		log.Debug("Update Status Base on Start Date: ", err)
		return err
	}

	logrus.Debug("Update At: ", &now)

	var lastID int64

	err = stmt.QueryRow(&now).Scan(&lastID)

	if err != nil {
		log.Debug("Update Status Base on Start Date: ", err)
		return err
	}

	return nil
}

func (m *psqlVoucherRepository) GetVoucher(c echo.Context, voucherID string) (*models.Voucher, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	result := new(models.Voucher)
	query := `SELECT c.id, c.name, c.description, c.start_date, c.end_date, c.point, c.value, c.image_url, c.stock, coalesce(d.available, 0), c.terms_and_conditions, c.how_to_use, c.limit_per_user
	FROM vouchers c LEFT JOIN(SELECT b.id, coalesce(count(a.id), 0) as available FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id 
	WHERE a.status = 0 GROUP BY b.id) d ON d.id = c.id WHERE c.id = $1`

	err := m.Conn.QueryRow(query, voucherID).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.StartDate,
		&result.EndDate,
		&result.Point,
		&result.Value,
		&result.ImageURL,
		&result.Stock,
		&result.Available,
		&result.TermsAndConditions,
		&result.HowToUse,
		&result.LimitPerUser,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return result, nil
}

func (m *psqlVoucherRepository) GetVouchersUser(c echo.Context, payload map[string]interface{}) ([]models.PromoCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	paging := ""
	where := ""
	query := `SELECT a.id, a.promo_code, a.bought_date, b.id, b.name, b.description, b.terms_and_conditions, b.how_to_use, b.limit_per_user, b.start_date, b.end_date, b.value, b.image_url 
	FROM promo_codes AS a LEFT JOIN vouchers AS b ON b.id = a.voucher_id WHERE a.promo_code IS NOT NULL AND a.status = 1`

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	if payload["userID"].(string) != "" {
		where += " AND user_id='" + payload["userID"].(string) + "'"
	}

	query += where + " ORDER BY a.bought_date DESC" + paging

	rows, err := m.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()
	var result []models.PromoCode

	for rows.Next() {
		var promoCode models.PromoCode
		var voucher models.Voucher
		err = rows.Scan(
			&promoCode.ID,
			&promoCode.PromoCode,
			&promoCode.BoughtDate,
			&voucher.ID,
			&voucher.Name,
			&voucher.Description,
			&voucher.TermsAndConditions,
			&voucher.HowToUse,
			&voucher.LimitPerUser,
			&voucher.StartDate,
			&voucher.EndDate,
			&voucher.Value,
			&voucher.ImageURL,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		if voucher.ID != 0 {
			promoCode.Voucher = &voucher
		}

		result = append(result, promoCode)
	}

	return result, nil
}

func (m *psqlVoucherRepository) CountVouchers(c echo.Context, payload map[string]interface{}, expired bool) (int, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	where := ""
	var total int
	query := `SELECT coalesce(COUNT(id), 0) FROM vouchers WHERE id IS NOT NULL AND status = 1`

	if payload["name"].(string) != "" {
		where += " AND name LIKE '%" + payload["name"].(string) + "%'"
	}

	if payload["startDate"].(string) != "" {
		where += " AND start_date::timestamp::date >= '" + payload["startDate"].(string) + "'"
	}

	if payload["endDate"].(string) != "" {
		where += " AND end_date::timestamp::date <= '" + payload["endDate"].(string) + "'"
	}

	if expired {
		where += " AND end_date::date >= now()"
	}

	query += where
	err := m.Conn.QueryRow(query).Scan(&total)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	return total, nil
}

func (m *psqlVoucherRepository) DeleteVoucher(c echo.Context, id int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `DELETE FROM vouchers WHERE ID = $1`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	result, err := stmt.Query(id)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	requestLogger.Debug("Result delete vouchers: ", result)

	return nil
}

func (m *psqlVoucherRepository) CountPromoCode(c echo.Context, payload map[string]interface{}) (int, error) {
	var total int
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	where := ""
	query := `SELECT coalesce(COUNT(id), 0) FROM promo_codes WHERE id IS NOT NULL`

	if payload["userID"].(string) != "" {
		where += " AND user_id='" + payload["userID"].(string) + "'"
	}

	if payload["status"].(string) != "" {
		where += " AND status='" + payload["status"].(string) + "'"
	}

	if payload["voucherID"].(string) != "" {
		where += " AND voucher_id='" + payload["voucherID"].(string) + "'"
	}

	query += where
	err := m.Conn.QueryRow(query).Scan(&total)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	return total, nil
}

func (m *psqlVoucherRepository) UpdatePromoCodeBought(c echo.Context, voucherID string, userID string) (*models.PromoCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	result := new(models.PromoCode)
	now := time.Now()

	querySelect := `SELECT id FROM promo_codes WHERE status = 0 AND voucher_id = $1 ORDER BY promo_code ASC LIMIT 1`
	err := m.Conn.QueryRow(querySelect, voucherID).Scan(&result.ID)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	queryUpdate := `UPDATE promo_codes SET status = 1, user_id = $1, bought_date = $2, updated_at = $3 WHERE id = $4 RETURNING promo_code, bought_date`
	stmt, err := m.Conn.Prepare(queryUpdate)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	err = stmt.QueryRow(userID, &now, &now, &result.ID).Scan(&result.PromoCode, &result.BoughtDate)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return result, nil
}

func (m *psqlVoucherRepository) UpdatePromoCodeRedeemed(c echo.Context, voucherID string, userID string, code string) (*models.PromoCode, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	result := new(models.PromoCode)
	queryUpdate := `UPDATE promo_codes SET status = 2, redeemed_date = $1, updated_at = $2 WHERE user_id = $3 AND promo_code = $4 RETURNING promo_code, redeemed_date`
	stmt, err := m.Conn.Prepare(queryUpdate)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	err = stmt.QueryRow(&now, &now, &userID, &code).Scan(&result.PromoCode, &result.RedeemedDate)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return result, nil
}

func (m *psqlVoucherRepository) GetVoucherCode(c echo.Context, voucherCode string, userID string) (*models.PromoCode, string, error) {
	var voucherID string
	result := &models.PromoCode{}
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT pc.id, pc.promo_code, pc.status, pc.user_id, pc.redeemed_date, pc.bought_date, pc.voucher_id
			  FROM promo_codes pc WHERE pc.promo_code = $1 AND pc.user_id = $2;`

	err := m.Conn.QueryRow(query, voucherCode, userID).Scan(
		&result.ID,
		&result.PromoCode,
		&result.Status,
		&result.UserID,
		&result.RedeemedDate,
		&result.BoughtDate,
		&voucherID,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, "", err
	}

	return result, voucherID, nil
}
