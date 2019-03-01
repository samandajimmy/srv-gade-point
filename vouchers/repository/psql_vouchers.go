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
	query := `SELECT d.id, d.name, d.description, d.start_date, d.end_date, d.point, d.journal_account, d.value, d.image_url, d.status, d.stock, d.prefix_promo_code, d.amount, CASE WHEN d.end_date::date < now()::date THEN 0 ELSE coalesce(e.avaliable,0) END AS avaliable, coalesce(f.bought,0) bought , coalesce(g.reedem,0) reedem, CASE WHEN coalesce(h.expired,0)-coalesce(g.reedem,0) < 0 THEN 0 ELSE coalesce(h.expired,0)-coalesce(g.reedem,0) END AS expired, d.validators, d.updated_at, d.created_at FROM (SELECT b.id, b.name, b.description, b.start_date, b.end_date, b.point, b.journal_account, b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, count(a.id) as amount, b.validators, b.updated_at, b.created_at FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id	GROUP BY b.id, b.name, b.description, b.start_date, b.end_date, b.point, b.journal_account, b.value, b.image_url, b.status, b.stock, b.prefix_promo_code, b.validators, b.updated_at, b.created_at) as d LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as avaliable FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 0	GROUP BY b.id) as e ON e.id = d.id LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as bought FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 1 GROUP BY b.id) as f ON f.id = d.id	LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as reedem FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE a.status = 2 GROUP BY b.id) as g ON g.id = d.id LEFT JOIN (SELECT b.id, coalesce(count(a.id), 0) as expired FROM promo_codes a LEFT JOIN vouchers b ON b.id=a.voucher_id WHERE end_date::date < now()::date GROUP BY b.id) as h ON h.id = d.id	WHERE d.id IS NOT NULL`

	where := ""

	if name != "" {
		where += " AND name LIKE '%" + name + "%'"
	}

	if status != "" {
		where += " AND status='" + status + "'"
	}

	if startDate != "" {
		where += " AND start_date >= '" + startDate + "'"
	}

	if endDate != "" {
		where += " AND end_date <= '" + endDate + "'"
	}

	query += where + " ORDER BY created_at DESC " + paging

	res, err := m.getVoucher(ctx, query)
	if err != nil {
		return nil, err
	}

	return res, err

}

// Execute query select from func GetVouchers return all data voucher with detail monitoring promo code
func (m *psqlVoucherRepository) getVoucher(ctx context.Context, query string) ([]*models.Voucher, error) {
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
			&t.Avaliable,
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
		valueArgs = append(valueArgs, promoCode.VoucherId)
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

	fmt.Println(result)

	return nil
}

// For count all data voucher by id
func (m *psqlVoucherRepository) CountVouchers(ctx context.Context, status string) (int, error) {

	query := `SELECT coalesce(COUNT(id), 0) FROM vouchers WHERE id IS NOT NULL`

	where := ""

	if status != "" {
		where += " AND status='" + status + "'"
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

	fmt.Println(result)

	return nil
}
