package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"time"

	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type psqlCampaignRepository struct {
	Conn *sql.DB
}

// NewPsqlCampaignRepository will create an object that represent the campaigns.Repository interface
func NewPsqlCampaignRepository(Conn *sql.DB) campaigns.Repository {
	return &psqlCampaignRepository{Conn}
}

func (m *psqlCampaignRepository) CreateCampaign(ctx context.Context, a *models.Campaign) error {

	query := `INSERT INTO campaigns (name, description, start_date, end_date, status, type, validators, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)  RETURNING id`
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

	err = stmt.QueryRowContext(ctx, a.Name, a.Description, a.StartDate, a.EndDate, a.Status, a.Type, string(validator), time.Now()).Scan(&lastID)
	if err != nil {
		return err
	}

	a.ID = lastID
	a.CreatedAt = &models.TimeNow
	return nil
}

func (m *psqlCampaignRepository) UpdateCampaign(ctx context.Context, id int64, updateCampaign *models.UpdateCampaign) error {

	query := `UPDATE campaigns SET status = $1, updated_at = $2 WHERE id = $3 RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	logrus.Debug("Update At: ", time.Now())

	var lastID int64

	err = stmt.QueryRowContext(ctx, updateCampaign.Status, time.Now(), id).Scan(&lastID)
	if err != nil {
		return err
	}

	return nil
}

func (m *psqlCampaignRepository) GetCampaign(ctx context.Context, name string, status string, startDate string, endDate string, page int, limit int) ([]*models.Campaign, error) {
	paging := ""
	where := ""
	query := `SELECT id, name, description, start_date, end_date, status, type, validators, updated_at, created_at FROM campaigns WHERE id IS NOT NULL`

	if page > 0 || limit > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", limit, ((page - 1) * limit))
	}

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

	query += where + " ORDER BY created_at DESC" + paging

	res, err := m.getCampaign(ctx, query)

	if err != nil {
		return nil, err
	}

	return res, err

}

func (m *psqlCampaignRepository) getCampaign(ctx context.Context, query string) ([]*models.Campaign, error) {
	var validator json.RawMessage
	result := make([]*models.Campaign, 0)
	rows, err := m.Conn.QueryContext(ctx, query)
	defer rows.Close()

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	for rows.Next() {
		t := new(models.Campaign)
		var createDate, updateDate pq.NullTime
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.StartDate,
			&t.EndDate,
			&t.Status,
			&t.Type,
			&validator,
			&updateDate,
			&createDate,
		)
		t.CreatedAt = &createDate.Time
		t.UpdatedAt = &updateDate.Time
		err = json.Unmarshal([]byte(validator), &t.Validators)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		result = append(result, t)
	}

	return result, nil
}

func (m *psqlCampaignRepository) GetValidatorCampaign(ctx context.Context, a *models.GetCampaignValue) (*models.Campaign, error) {
	result := new(models.Campaign)
	var validator json.RawMessage

	query := `SELECT id, validators FROM campaigns WHERE status = 1 AND start_date::date <= now()::date
	AND end_date::date >= now()::date AND validators->>'channel'=$1 AND validators->>'product'=$2 AND validators->>'transactionType'=$3 AND validators->>'unit'=$4 ORDER BY end_date ASC LIMIT 1`

	err := m.Conn.QueryRowContext(ctx, query, a.Channel, a.Product, a.TransactionType, a.Unit).Scan(&result.ID, &validator)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(validator), &result.Validators)
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (m *psqlCampaignRepository) SavePoint(ctx context.Context, a *models.CampaignTrx) error {
	query := `INSERT INTO campaign_transactions (user_id, point_amount, transaction_type, transaction_date, campaign_id, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)  RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	logrus.Debug("Created At: ", time.Now())

	var lastID int64

	err = stmt.QueryRowContext(ctx, a.UserID, a.PointAmount, a.TransactionType, a.TransactionDate, a.Campaign.ID, a.CreatedAt).Scan(&lastID)
	if err != nil {
		return err
	}

	a.ID = lastID
	return nil
}

func (m *psqlCampaignRepository) GetUserPoint(ctx context.Context, UserId string) (float64, error) {
	var pointDebet float64
	var pointKredit float64

	queryDebet := `SELECT coalesce(sum(point_amount), 0) as debet FROM public.campaign_transactions WHERE user_id = $1 AND transaction_type = 'D' AND to_char(transaction_date, 'YYYY') = to_char(NOW(), 'YYYY')`

	err := m.Conn.QueryRowContext(ctx, queryDebet, UserId).Scan(&pointDebet)
	if err != nil {
		return 0, err
	}

	queryKredit := `SELECT coalesce(sum(point_amount), 0) as debet FROM public.campaign_transactions WHERE user_id = $1 AND transaction_type = 'K' AND to_char(transaction_date, 'YYYY') = to_char(NOW(), 'YYYY')`

	err = m.Conn.QueryRowContext(ctx, queryKredit, UserId).Scan(&pointKredit)
	if err != nil {
		return 0, err
	}

	pointAmount := pointDebet - pointKredit

	return pointAmount, nil
}

func (m *psqlCampaignRepository) GetUserPointHistory(ctx context.Context, userID string) ([]models.CampaignTrx, error) {
	var dataHistory []models.CampaignTrx

	query := `select
				ct.id,
				ct.user_id,
				ct.point_amount,
				ct.transaction_type,
				ct.transaction_date,
				coalesce(ct.campaign_id, 0) campaign_id,
				coalesce(c.name, '') campaign_name,
				coalesce(c.description, '') campaign_description,
				coalesce(ct.promo_code_id, 0) promo_code_id,
				coalesce(pc.promo_code, '') promo_code,
				coalesce(v.name, '') voucher_name,
				coalesce(v.description, '') voucher_description
			from
				campaign_transactions ct
			left join campaigns c on
				ct.campaign_id = c.id
			left join promo_codes pc on
				pc.id = ct.promo_code_id
			left join vouchers v on
				pc.voucher_id = v.id
			where
				ct.user_id = $1
			order by
				ct.transaction_date desc;`

	rows, err := m.Conn.QueryContext(ctx, query, userID)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	for rows.Next() {
		var ct models.CampaignTrx
		var campaign models.Campaign
		var promoCodes models.PromoCode
		var voucher models.Voucher

		err = rows.Scan(
			&ct.ID,
			&ct.UserID,
			&ct.PointAmount,
			&ct.TransactionType,
			&ct.TransactionDate,
			&campaign.ID,
			&campaign.Name,
			&campaign.Description,
			&promoCodes.ID,
			&promoCodes.PromoCode,
			&voucher.Name,
			&voucher.Description,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}

		if campaign.ID != 0 {
			ct.Campaign = &campaign
		}

		if promoCodes.ID != 0 {
			ct.PromoCode = &promoCodes
		}

		dataHistory = append(dataHistory, ct)
	}

	return dataHistory, nil
}

// Count data campaign
func (m *psqlCampaignRepository) CountCampaign(ctx context.Context, name string, status string, startDate string, endDate string) (int, error) {

	query := `SELECT coalesce(COUNT(id), 0) FROM campaigns WHERE id IS NOT NULL`

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

	query += where

	var total int

	err := m.Conn.QueryRowContext(ctx, query).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}
