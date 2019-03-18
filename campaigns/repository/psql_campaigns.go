package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"time"

	"github.com/labstack/gommon/log"
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
	now := time.Now()
	query := `INSERT INTO campaigns (name, description, start_date, end_date, status, type, validators, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)  RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	logrus.Debug("Created At: ", &now)

	var lastID int64
	validator, err := json.Marshal(a.Validators)

	if err != nil {
		return err
	}

	err = stmt.QueryRowContext(ctx, a.Name, a.Description, a.StartDate, a.EndDate, a.Status, a.Type, string(validator), &now).Scan(&lastID)

	if err != nil {
		return err
	}

	a.ID = lastID
	a.CreatedAt = &now
	return nil
}

func (m *psqlCampaignRepository) UpdateCampaign(ctx context.Context, id int64, updateCampaign *models.UpdateCampaign) error {
	now := time.Now()
	query := `UPDATE campaigns SET status = $1, updated_at = $2 WHERE id = $3 RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	logrus.Debug("Update At: ", &now)

	var lastID int64

	err = stmt.QueryRowContext(ctx, updateCampaign.Status, &now, id).Scan(&lastID)

	if err != nil {
		return err
	}

	return nil
}

func (m *psqlCampaignRepository) UpdateExpiryDate(ctx context.Context) error {
	now := time.Now()
	query := `UPDATE campaigns SET status = 0, updated_at = $1 WHERE end_date::timestamp::date < now()::date`
	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		log.Debug("Update Status Base on Expiry Date: ", err)
		return err
	}

	logrus.Debug("Update At: ", &now)

	var lastID int64

	err = stmt.QueryRowContext(ctx, &now).Scan(&lastID)

	if err != nil {
		log.Debug("Update Status Base on Expiry Date: ", err)
		return err
	}

	return nil
}

func (m *psqlCampaignRepository) UpdateStatusBasedOnStartDate() error {
	now := time.Now()
	query := `UPDATE campaigns SET status = 1, updated_at = $1 WHERE start_date::timestamp::date = now()::date`
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
	var validator json.RawMessage
	result := new(models.Campaign)
	reqValidator := &models.Validator{}
	query := `SELECT id, validators FROM campaigns WHERE status = 1 AND start_date::date <= now()::date AND end_date::date >= now()::date`

	mapValidator := reqValidator.GetValidatorKeys(a)

	for key, value := range mapValidator {
		query += ` AND validators->>'` + key + `'='` + value + `'`
	}

	query += ` ORDER BY end_date ASC LIMIT 1`
	err := m.Conn.QueryRowContext(ctx, query).Scan(&result.ID, &validator)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	err = json.Unmarshal([]byte(validator), &result.Validators)

	if err != nil {
		log.Error(err)
		return nil, err
	}

	return result, nil

}

func (m *psqlCampaignRepository) SavePoint(ctx context.Context, cmpgnTrx *models.CampaignTrx) error {
	now := time.Now()
	var id int64
	var query string

	if cmpgnTrx.TransactionType == models.TransactionPointTypeDebet {
		query = `INSERT INTO campaign_transactions (user_id, point_amount, transaction_type, transaction_date, campaign_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)  RETURNING id`
		id = cmpgnTrx.Campaign.ID
	}

	if cmpgnTrx.TransactionType == models.TransactionPointTypeKredit {
		query = `INSERT INTO campaign_transactions (user_id, point_amount, transaction_type, transaction_date, promo_code_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)  RETURNING id`
		id = cmpgnTrx.PromoCode.ID
	}

	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	logrus.Debug("Created At: ", &now)

	var lastID int64
	err = stmt.QueryRowContext(ctx, cmpgnTrx.UserID, cmpgnTrx.PointAmount, cmpgnTrx.TransactionType, cmpgnTrx.TransactionDate, id, cmpgnTrx.CreatedAt).Scan(&lastID)

	if err != nil {
		return err
	}

	cmpgnTrx.ID = lastID
	return nil
}

func (m *psqlCampaignRepository) GetUserPoint(ctx context.Context, UserID string) (float64, error) {
	var pointDebet float64
	var pointKredit float64
	queryDebet := `SELECT coalesce(sum(point_amount), 0) as debet FROM public.campaign_transactions WHERE user_id = $1 AND transaction_type = 'D' AND to_char(transaction_date, 'YYYY') = to_char(NOW(), 'YYYY')`
	err := m.Conn.QueryRowContext(ctx, queryDebet, UserID).Scan(&pointDebet)

	if err != nil {
		return 0, err
	}

	queryKredit := `SELECT coalesce(sum(point_amount), 0) as debet FROM public.campaign_transactions WHERE user_id = $1 AND transaction_type = 'K' AND to_char(transaction_date, 'YYYY') = to_char(NOW(), 'YYYY')`
	err = m.Conn.QueryRowContext(ctx, queryKredit, UserID).Scan(&pointKredit)

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
				coalesce(pc.voucher_id, 0) voucher_id,
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
			&voucher.ID,
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
			ct.PromoCode.ID = 0 // remove promo codes ID from the response
		}

		if voucher.ID != 0 {
			ct.PromoCode.Voucher = &voucher
		}

		dataHistory = append(dataHistory, ct)
	}

	return dataHistory, nil
}

func (m *psqlCampaignRepository) CountCampaign(ctx context.Context, name string, status string, startDate string, endDate string) (int, error) {
	var total int
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
	err := m.Conn.QueryRowContext(ctx, query).Scan(&total)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (m *psqlCampaignRepository) GetCampaignDetail(ctx context.Context, id int64) (*models.Campaign, error) {
	var validator json.RawMessage
	var createDate, updateDate pq.NullTime
	result := new(models.Campaign)

	query := `SELECT id, name, description, start_date, end_date, status, type, validators, updated_at, created_at FROM campaigns WHERE id = $1`

	err := m.Conn.QueryRowContext(ctx, query, id).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.StartDate,
		&result.EndDate,
		&result.Status,
		&result.Type,
		&validator,
		&createDate,
		&updateDate,
	)

	result.CreatedAt = &createDate.Time
	result.UpdatedAt = &updateDate.Time
	err = json.Unmarshal([]byte(validator), &result.Validators)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	return result, err
}
