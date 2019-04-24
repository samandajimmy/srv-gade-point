package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
	"time"

	"github.com/labstack/echo"
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

func (m *psqlCampaignRepository) CreateCampaign(c echo.Context, campaign *models.Campaign) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `INSERT INTO campaigns (name, description, start_date, end_date, status, type, validators, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	var lastID int64
	validator, err := json.Marshal(campaign.Validators)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(campaign.Name, campaign.Description, campaign.StartDate, campaign.EndDate, campaign.Status, campaign.Type, string(validator), &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	campaign.ID = lastID
	campaign.CreatedAt = &now
	return nil
}

func (m *psqlCampaignRepository) UpdateCampaign(c echo.Context, id int64, updateCampaign *models.UpdateCampaign) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `UPDATE campaigns SET status = $1, updated_at = $2 WHERE id = $3 RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	var lastID int64
	err = stmt.QueryRow(updateCampaign.Status, &now, id).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func (m *psqlCampaignRepository) UpdateExpiryDate(c echo.Context) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `UPDATE campaigns SET status = 0, updated_at = $1 WHERE end_date::timestamp::date < now()::date AND status = 1`
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

func (m *psqlCampaignRepository) UpdateStatusBasedOnStartDate() error {
	now := time.Now()
	query := `UPDATE campaigns SET status = 1, updated_at = $1 WHERE start_date::timestamp::date = now()::date`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		logrus.Debug("Update Status Base on Start Date: ", err)

		return err
	}

	var lastID int64
	err = stmt.QueryRow(&now).Scan(&lastID)

	if err != nil {
		logrus.Debug("Update Status Base on Start Date: ", err)

		return err
	}

	return nil
}

func (m *psqlCampaignRepository) GetCampaign(c echo.Context, payload map[string]interface{}) ([]*models.Campaign, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	paging := ""
	where := ""
	query := `SELECT id, name, description, start_date, end_date, status, type, validators, updated_at, created_at FROM campaigns WHERE id IS NOT NULL`

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	if payload["name"].(string) != "" {
		where += " AND name LIKE '%" + payload["name"].(string) + "%'"
	}

	if payload["status"].(string) != "" {
		where += " AND status='" + payload["status"].(string) + "'"
	}

	if payload["startDate"].(string) != "" {
		where += " AND start_date >= '" + payload["startDate"].(string) + "'"
	}

	if payload["endDate"].(string) != "" {
		where += " AND end_date <= '" + payload["endDate"].(string) + "'"
	}

	query += where + " ORDER BY created_at DESC" + paging
	res, err := m.getCampaign(c, query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return res, err

}

func (m *psqlCampaignRepository) getCampaign(c echo.Context, query string) ([]*models.Campaign, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var validator json.RawMessage
	result := make([]*models.Campaign, 0)
	rows, err := m.Conn.Query(query)
	defer rows.Close()

	if err != nil {
		requestLogger.Debug(err)

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
			requestLogger.Debug(err)

			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (m *psqlCampaignRepository) GetValidatorCampaign(c echo.Context, payload *models.GetCampaignValue) (*models.Campaign, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var validator json.RawMessage
	result := new(models.Campaign)
	query := `SELECT id, validators FROM campaigns WHERE status = 1 AND start_date::date <= now()::date
	AND end_date::date >= now()::date AND validators->>'channel'=$1 AND validators->>'product'=$2 AND validators->>'transactionType'=$3 AND validators->>'unit'=$4 ORDER BY start_date DESC LIMIT 1`
	err := m.Conn.QueryRow(query, payload.Channel, payload.Product, payload.TransactionType, payload.Unit).Scan(&result.ID, &validator)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	err = json.Unmarshal([]byte(validator), &result.Validators)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return result, nil
}

func (m *psqlCampaignRepository) SavePoint(c echo.Context, cmpgnTrx *models.CampaignTrx) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, cmpgnTrx)
	now := time.Now()
	var id int64
	var query string

	if cmpgnTrx.TransactionType == models.TransactionPointTypeDebet {
		query = `INSERT INTO campaign_transactions (user_id, point_amount, transaction_type, transaction_date, reff_core, campaign_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING id`
		id = cmpgnTrx.Campaign.ID
	}

	if cmpgnTrx.TransactionType == models.TransactionPointTypeKredit {
		query = `INSERT INTO campaign_transactions (user_id, point_amount, transaction_type, transaction_date, reff_core, promo_code_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING id`
		id = cmpgnTrx.PromoCode.ID
	}

	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	cmpgnTrx.CreatedAt = &now
	var lastID int64
	err = stmt.QueryRow(cmpgnTrx.UserID, cmpgnTrx.PointAmount, cmpgnTrx.TransactionType, cmpgnTrx.TransactionDate, cmpgnTrx.ReffCore, id, cmpgnTrx.CreatedAt).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	cmpgnTrx.ID = lastID
	return nil
}

func (m *psqlCampaignRepository) GetUserPoint(c echo.Context, UserID string) (float64, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var pointDebet float64
	var pointKredit float64
	queryDebet := `SELECT coalesce(sum(point_amount), 0) as debet FROM public.campaign_transactions WHERE user_id = $1 AND transaction_type = 'D' AND to_char(transaction_date, 'YYYY') = to_char(NOW(), 'YYYY')`
	err := m.Conn.QueryRow(queryDebet, UserID).Scan(&pointDebet)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	queryKredit := `SELECT coalesce(sum(point_amount), 0) as debet FROM public.campaign_transactions WHERE user_id = $1 AND transaction_type = 'K' AND to_char(transaction_date, 'YYYY') = to_char(NOW(), 'YYYY')`
	err = m.Conn.QueryRow(queryKredit, UserID).Scan(&pointKredit)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	pointAmount := pointDebet - pointKredit
	return pointAmount, nil
}

func (m *psqlCampaignRepository) GetUserPointHistory(c echo.Context, payload map[string]interface{}) ([]models.CampaignTrx, error) {
	var dataHistory []models.CampaignTrx
	where := ""
	paging := ""
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	startDateRg := payload["startDateRg"].(string)
	endDateRg := payload["endDateRg"].(string)

	query := `select
				ct.id,
				ct.user_id,
				ct.point_amount,
				ct.transaction_type,
				ct.transaction_date,
				coalesce(ct.reff_core, '') reff_core,
				coalesce(ct.campaign_id, 0) campaign_id,
				coalesce(c.name, '') campaign_name,
				coalesce(c.description, '') campaign_description,
				coalesce(ct.promo_code_id, 0) promo_code_id,
				coalesce(pc.promo_code, '') promo_code,
				coalesce(pc.voucher_id, 0) voucher_id,
				coalesce(v.name, '') voucher_name,
				coalesce(v.description, '') voucher_description
			from campaign_transactions ct
			left join campaigns c on ct.campaign_id = c.id
			left join promo_codes pc on pc.id = ct.promo_code_id
			left join vouchers v on pc.voucher_id = v.id
			where ct.user_id = $1`

	if startDateRg != "" {
		where += " and ct.transaction_date::timestamp::date >= '" + startDateRg + "'"
	}

	if endDateRg != "" {
		where += " and ct.transaction_date::timestamp::date <= '" + endDateRg + "'"
	}

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
	}

	query += where + " order by ct.transaction_date desc" + paging + ";"
	rows, err := m.Conn.Query(query, payload["userID"].(string))

	if err != nil {
		requestLogger.Debug(err)

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
			&ct.ReffCore,
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
			requestLogger.Debug(err)

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

func (m *psqlCampaignRepository) CountUserPointHistory(c echo.Context, payload map[string]interface{}) (string, error) {
	var counter string
	where := ""
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `select COUNT(*) counter from campaign_transactions where user_id = $1`

	if payload["startDateRg"].(string) != "" {
		where += " and transaction_date::timestamp::date >= '" + payload["startDateRg"].(string) + "'"
	}

	if payload["endDateRg"].(string) != "" {
		where += " and transaction_date::timestamp::date <= '" + payload["endDateRg"].(string) + "'"
	}

	query += where + ";"
	err := m.Conn.QueryRow(query, payload["userID"].(string)).Scan(&counter)

	if err != nil {
		requestLogger.Debug(err)

		return "", err
	}

	return counter, nil
}

func (m *psqlCampaignRepository) CountCampaign(c echo.Context, payload map[string]interface{}) (int, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var total int
	query := `SELECT coalesce(COUNT(id), 0) FROM campaigns WHERE id IS NOT NULL`
	where := ""

	if payload["name"].(string) != "" {
		where += " AND name LIKE '%" + payload["name"].(string) + "%'"
	}

	if payload["status"].(string) != "" {
		where += " AND status='" + payload["status"].(string) + "'"
	}

	if payload["startDate"].(string) != "" {
		where += " AND start_date >= '" + payload["startDate"].(string) + "'"
	}

	if payload["endDate"].(string) != "" {
		where += " AND end_date <= '" + payload["endDate"].(string) + "'"
	}

	query += where
	err := m.Conn.QueryRow(query).Scan(&total)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	return total, nil
}

func (m *psqlCampaignRepository) GetCampaignDetail(c echo.Context, id int64) (*models.Campaign, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var validator json.RawMessage
	var createDate, updateDate pq.NullTime
	result := new(models.Campaign)

	query := `SELECT id, name, description, start_date, end_date, status, type, validators, updated_at, created_at FROM campaigns WHERE id = $1`

	err := m.Conn.QueryRow(query, id).Scan(
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

	return result, nil
}
