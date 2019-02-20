package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/vouchers"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type psqlVoucherRepository struct {
	Conn *sql.DB
}

type JSONB []byte

// NewPsqlVoucherRepository will create an object that represent the vouchers. Repository interface
func NewPsqlVoucherRepository(Conn *sql.DB) vouchers.Repository {
	return &psqlVoucherRepository{Conn}
}

func (m *psqlVoucherRepository) CreateVoucher(ctx context.Context, a *models.Voucher) error {

	query := `INSERT INTO vouchers (name, description, start_date, end_date, point, journal_account, value, image_url, status, validators, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)  RETURNING id`
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

	err = stmt.QueryRowContext(ctx, a.Name, a.Description, a.StartDate, a.EndDate, a.Point, a.JournalAccount, a.Value, a.ImageUrl, a.Status, string(validator), time.Now()).Scan(&lastID)
	if err != nil {
		return err
	}

	a.ID = lastID
	a.CreatedAt = time.Now()
	return nil
}

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

func (m *psqlVoucherRepository) GetVoucher(ctx context.Context, name string, status string, startDate string, endDate string) ([]*models.Voucher, error) {
	query := `SELECT * FROM vouchers WHERE id IS NOT NULL`

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

	query += where + " ORDER BY created_at DESC"

	res, err := m.getVoucher(ctx, query)
	if err != nil {
		return nil, err
	}

	return res, err

}

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
