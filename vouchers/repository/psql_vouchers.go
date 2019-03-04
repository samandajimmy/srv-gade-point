package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat           = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
	fieldInsertPromoCode = 4
)

var (
	count = []int{1, 2, 3, 4}
)

type psqlVoucherRepository struct {
	Conn *sql.DB
}

type JSONB []byte

// NewPsqlVoucherRepository will create an object that represent the vouchers. Repository interface
func NewPsqlVoucherRepository(Conn *sql.DB) vouchers.Repository {
	return &psqlVoucherRepository{Conn}
}

// Insert new voucher to database table vouchers
func (m *psqlVoucherRepository) CreateVoucher(ctx context.Context, a *models.Voucher) error {

	query := `INSERT INTO vouchers (name, description, start_date, end_date, point, journal_account, value, image_url, status, stock, prefix_promo_code, validators, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)  RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	logrus.Debug("Created At: ", time.Now())

	var lastID int64
	validator, err := json.Marshal(a.Validators)
	if err != nil {
		return err
	}

	err = stmt.QueryRowContext(ctx, a.Name, a.Description, a.StartDate, a.EndDate, a.Point, a.JournalAccount, a.Value, a.ImageUrl, a.Status, a.Stock, a.PrefixPromoCode, string(validator), time.Now()).Scan(&lastID)
	if err != nil {
		return err
	}

	a.ID = lastID
	a.CreatedAt = time.Now()
	return nil
}

// Update status voucher to database table vouchers
func (m *psqlVoucherRepository) UpdateVoucher(ctx context.Context, id int64, updateVoucher *models.UpdateVoucher) error {

	query := `UPDATE vouchers SET status = $1, updated_at = $2 WHERE id = $3 RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	logrus.Debug("Update At: ", time.Now())

	var lastID int64

	err = stmt.QueryRowContext(ctx, updateVoucher.Status, time.Now(), id).Scan(&lastID)
	if err != nil {
		return err
	}

	return nil
}

// Get data all voucher with detail monitoring promo code
func (m *psqlVoucherRepository) GetVouchers(ctx context.Context, name string, status string, startDate string, endDate string, page int32, limit int32) ([]*models.Voucher, error) {

	paging := fmt.Sprintf(" LIMIT %d OFFSET %d", limit, ((page - 1) * limit))
	query := `SELECT d.id, d.name, d.description, d.start_date, d.end_date, d.point, d.journal_account, d.value, d.image_url, d.status, d.stock, d.prefix_promo_code, d.amount, CASE WHEN d.end_date::date < now()::date THEN 0 ELSE coalesce(e.available,0) END AS available, coalesce(f.bought,0) bought , coalesce(g.reedem,0) reedem, CASE WHEN coalesce(h.expired,0)-coalesce(g.reedem,0) < 0 THEN 0 ELSE coalesce(h.expired,0)-coalesce(g.reedem,0) END AS expired, d.validators, d.updated_at, d.created_at FROM (SELECT b.id, b.name, b.description, b.start_date, b.end_date, b.point, b.journal_account, b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, count(a.id) as amount, b.validators, b.updated_at, b.created_at FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id	GROUP BY b.id, b.name, b.description, b.start_date, b.end_date, b.point, b.journal_account, b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, b.validators, b.updated_at, b.created_at) as d LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as available FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 0	GROUP BY b.id) as e ON e.id = d.id LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as bought FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 1 GROUP BY b.id) as f ON f.id = d.id	LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as reedem FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 2 GROUP BY b.id) as g ON g.id = d.id LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as expired FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE end_date::date < now()::date GROUP BY b.id) as h ON h.id = d.id	WHERE d.id IS NOT NULL`

	where := ""

	if name != "" {
		where += " AND d.name LIKE '%" + name + "%'"
	}

	if status != "" {
		where += " AND d.status='" + status + "'"
	}

	if startDate != "" {
		where += " AND d.start_date >= '" + startDate + "'"
	}

	if endDate != "" {
		where += " AND d.end_date <= '" + endDate + "'"
	}

	query += where + " ORDER BY d.created_at DESC " + paging

	res, err := m.getVouchers(ctx, query)
	if err != nil {
		return nil, err
	}

	return res, err

}

// Execute query select from func GetVouchers return all data voucher with detail monitoring promo code
func (m *psqlVoucherRepository) getVouchers(ctx context.Context, query string) ([]*models.Voucher, error) {
	var validator json.RawMessage
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
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
			&t.ImageUrl,
			&t.Status,
			&t.Stock,
			&t.PrefixPromoCode,
			&t.Amount,
			&t.Available,
			&t.Bought,
			&t.Redeemed,
			&t.Expired,
			&validator,
			&updateDate,
			&createDate,
		)
		t.CreatedAt = createDate.Time
		t.UpdatedAt = updateDate.Time
		err = json.Unmarshal([]byte(validator), &t.Validators)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

// Insert data promo code from result generate promo code
func (m *psqlVoucherRepository) CreatePromoCode(ctx context.Context, promoCodes []*models.PromoCode) error {

	var valueStrings []string
	var valueArgs []interface{}
	i := 0
	for _, promoCode := range promoCodes {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*fieldInsertPromoCode+count[0], i*fieldInsertPromoCode+count[1], i*fieldInsertPromoCode+count[2], i*fieldInsertPromoCode+count[3]))
		valueArgs = append(valueArgs, promoCode.PromoCode)
		valueArgs = append(valueArgs, promoCode.Status)
		valueArgs = append(valueArgs, promoCode.Voucher.ID)
		valueArgs = append(valueArgs, promoCode.CreatedAt)
		i++
	}

	query := fmt.Sprintf("INSERT INTO promo_codes (promo_code, status, voucher_id, created_at) VALUES %s", strings.Join(valueStrings, ","))

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	result, err := stmt.QueryContext(ctx, valueArgs...)
	if err != nil {
		return err
	}
	logrus.Debug("Result created promo code: ", result)

	return nil
}

// For count all data voucher by id
func (m *psqlVoucherRepository) CountVouchers(ctx context.Context, status string, expired bool) (int, error) {

	query := `SELECT coalesce(COUNT(id), 0) FROM vouchers WHERE id IS NOT NULL`

	where := ""

	if status != "" {
		where += " AND status='" + status + "'"
	}
	if expired {
		where += " AND end_date::date >= now()"

	}

	query += where

	var total int

	err := m.Conn.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Delete voucher when failed generate promocode
func (m *psqlVoucherRepository) DeleteVoucher(ctx context.Context, id int64) error {

	query := `DELETE FROM vouchers WHERE ID = $1`

	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	result, err := stmt.QueryContext(ctx, id)
	if err != nil {
		return err
	}

	logrus.Debug("Result delete vouchers: ", result)

	return nil
}

// Get data all voucher for external
func (m *psqlVoucherRepository) GetVouchersExternal(ctx context.Context, name string, startDate string, endDate string, page int32, limit int32) ([]*models.VoucherDetail, error) {

	paging := fmt.Sprintf(" LIMIT %d OFFSET %d", limit, ((page - 1) * limit))
	query := `SELECT c.id, c.name, c.description, c.start_date, c.end_date, c.point, c.value, c.image_url, c.stock, d. available FROM vouchers c 
	LEFT JOIN(SELECT b.id, coalesce(count(a.id), 0) as available FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 0 GROUP BY b.id) d ON d.id = c.id WHERE c.status = 1 AND c.end_date::date >= now() `

	where := ""

	if name != "" {
		where += " AND name LIKE '%" + name + "%'"
	}

	if startDate != "" {
		where += " AND start_date >= '" + startDate + "'"
	}

	if endDate != "" {
		where += " AND end_date <= '" + endDate + "'"
	}

	query += where + " ORDER BY c.created_at DESC " + paging

	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.VoucherDetail, 0)

	for rows.Next() {
		t := new(models.VoucherDetail)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.StartDate,
			&t.EndDate,
			&t.Point,
			&t.Value,
			&t.ImageUrl,
			&t.Stock,
			&t.Available,
		)

		result = append(result, t)
	}

	return result, nil
}

// Get data voucher detail for external
func (m *psqlVoucherRepository) GetVoucher(ctx context.Context, voucherId string) (*models.Voucher, error) {
	result := new(models.Voucher)
	var validator json.RawMessage
	var createDate, updateDate pq.NullTime

	query := `SELECT d.id, d.name, d.description, d.start_date, d.end_date, d.point, d.journal_account, d.value, d.image_url, d.status, d.stock, d.prefix_promo_code, d.amount, CASE WHEN d.end_date::date < now()::date THEN 0 ELSE coalesce(e.available,0) END AS available, coalesce(f.bought,0) bought , coalesce(g.reedem,0) reedem, CASE WHEN coalesce(h.expired,0)-coalesce(g.reedem,0) < 0 THEN 0 ELSE coalesce(h.expired,0)-coalesce(g.reedem,0) END AS expired, d.validators, d.updated_at, d.created_at FROM (SELECT b.id, b.name, b.description, b.start_date, b.end_date, b.point, b.journal_account, b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, count(a.id) as amount, b.validators, b.updated_at, b.created_at FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id	GROUP BY b.id, b.name, b.description, b.start_date, b.end_date, b.point, b.journal_account, b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, b.validators, b.updated_at, b.created_at) as d LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as available FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 0	GROUP BY b.id) as e ON e.id = d.id LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as bought FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 1 GROUP BY b.id) as f ON f.id = d.id	LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as reedem FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 2 GROUP BY b.id) as g ON g.id = d.id LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as expired FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE end_date::date < now()::date GROUP BY b.id) as h ON h.id = d.id	WHERE d.id = $1`

	err := m.Conn.QueryRowContext(ctx, query, voucherId).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.StartDate,
		&result.EndDate,
		&result.Point,
		&result.JournalAccount,
		&result.Value,
		&result.ImageUrl,
		&result.Status,
		&result.Stock,
		&result.PrefixPromoCode,
		&result.Amount,
		&result.Available,
		&result.Bought,
		&result.Redeemed,
		&result.Expired,
		&validator,
		&updateDate,
		&createDate)
	if err != nil {
		return nil, err
	}

	result.CreatedAt = createDate.Time
	result.UpdatedAt = updateDate.Time
	err = json.Unmarshal([]byte(validator), &result.Validators)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return result, err

}

// Get data voucher detail for admin
func (m *psqlVoucherRepository) GetVoucherExternal(ctx context.Context, voucherId string) (*models.VoucherDetail, error) {
	result := new(models.VoucherDetail)
	query := `SELECT c.id, c.name, c.description, c.start_date, c.end_date, c.point, c.value, c.image_url, c.stock, d. available FROM vouchers c 
	LEFT JOIN(SELECT b.id, coalesce(count(a.id), 0) as available FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 0 GROUP BY b.id) d ON d.id = c.id WHERE c.id = $1 `

	err := m.Conn.QueryRowContext(ctx, query, voucherId).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.StartDate,
		&result.EndDate,
		&result.Point,
		&result.Value,
		&result.ImageUrl,
		&result.Stock,
		&result.Available,
	)
	if err != nil {
		return nil, err
	}

	return result, err
}

//Get vouchers user
func (m *psqlVoucherRepository) GetVouchersUser(ctx context.Context, userId string, status string, page int32, limit int32) ([]*models.VoucherUser, error) {
	paging := fmt.Sprintf(" LIMIT %d OFFSET %d", limit, ((page - 1) * limit))

	query := `SELECT a.promo_code, a.bought_date, b.name, b.description, b.start_date, b.end_date, b.value, b.image_url FROM promo_codes AS a LEFT JOIN vouchers AS b ON b.id = a.voucher_id WHERE a.promo_code IS NOT NULL`
	where := ""
	if status != "" {
		where += " AND a.status='" + status + "'"
	}

	if userId != "" {
		where += " AND a.user_id='" + userId + "'"
	}

	query += where + " ORDER BY a.bought_date DESC" + paging

	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.VoucherUser, 0)

	for rows.Next() {
		t := new(models.VoucherUser)
		err = rows.Scan(
			&t.PromoCode,
			&t.BoughtDate,
			&t.Name,
			&t.Description,
			&t.StartDate,
			&t.EndDate,
			&t.Value,
			&t.ImageUrl,
		)

		result = append(result, t)
	}

	return result, nil
}

// Count data promo code by userand status
func (m *psqlVoucherRepository) CountPromoCode(ctx context.Context, status string, userId string) (int, error) {

	query := `SELECT coalesce(COUNT(id), 0) FROM promo_codes WHERE id IS NOT NULL`

	where := ""

	if status != "" {
		where += " AND status='" + status + "'"
	}

	if userId != "" {
		where += " AND user_id='" + userId + "'"
	}

	query += where

	var total int

	err := m.Conn.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// Update promo code foruser
func (m *psqlVoucherRepository) UpdatePromoCodeBought(ctx context.Context, voucherId string, userId string) (int64, string, string, error) {
	var idPromoCode int64
	var promoCode string
	var boughtDate string

	querySelect := `SELECT id FROM promo_codes WHERE status = 0 AND voucher_id = $1 ORDER BY promo_code ASC LIMIT 1`

	err := m.Conn.QueryRowContext(ctx, querySelect, voucherId).Scan(&idPromoCode)
	if err != nil {
		return 0, "", "", err
	}

	queryUpdate := `UPDATE promo_codes SET status = 1, user_id = $1, bought_date = $2, updated_at = $3 WHERE id = $4 RETURNING promo_code, bought_date`
	stmt, err := m.Conn.PrepareContext(ctx, queryUpdate)
	if err != nil {
		return 0, "", "", err
	}

	logrus.Debug("Update At promo_codes : ", time.Now())

	err = stmt.QueryRowContext(ctx, userId, time.Now(), time.Now(), idPromoCode).Scan(&promoCode, &boughtDate)
	if err != nil {
		return 0, "", "", err
	}

	return idPromoCode, promoCode, boughtDate, nil
}
