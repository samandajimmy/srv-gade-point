package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/database"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewards"
	"strings"
	"time"

	"github.com/labstack/echo"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

type psqlCampaignRepository struct {
	Conn    *sql.DB
	dbBun   *database.DbBun
	rwdRepo rewards.RRepository
}

// NewPsqlCampaignRepository will create an object that represent the campaigns.Repository interface
func NewPsqlCampaignRepository(Conn *sql.DB, dbBun *database.DbBun, rwdRepo rewards.RRepository) campaigns.CRepository {
	return &psqlCampaignRepository{Conn, dbBun, rwdRepo}
}

func (m *psqlCampaignRepository) CreateCampaign(c echo.Context, campaign *models.Campaign) error {
	var endDate *string
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `INSERT INTO campaigns (name, description, start_date, end_date, status, created_at, metadata)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	var lastID int64

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	if campaign.EndDate != "" {
		endDate = &campaign.EndDate
	}

	metadata, err := json.Marshal(campaign.Metadata)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(campaign.Name, campaign.Description, campaign.StartDate, endDate,
		campaign.Status, &now, string(metadata)).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	campaign.ID = lastID
	campaign.CreatedAt = &now
	return nil
}

func (m *psqlCampaignRepository) UpdateCampaign(c echo.Context, id int64,
	updateCampaign *models.Campaign) error {

	var lastID int64
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `UPDATE campaigns SET status = $1, updated_at = $2 WHERE id = $3 RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(updateCampaign.Status, &now, id).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func (m *psqlCampaignRepository) UpdateExpiryDate(c echo.Context) error {
	var lastID int64
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `UPDATE campaigns SET status = 0, updated_at = $1
		WHERE end_date::timestamp::date < now()::date AND status = 1`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug("Update Status Base on Expiry Date: ", err)

		return err
	}

	err = stmt.QueryRow(&now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug("Update Status Base on Expiry Date: ", err)

		return err
	}

	return nil
}

func (m *psqlCampaignRepository) UpdateStatusBasedOnStartDate() error {
	var lastID int64
	now := time.Now()
	query := `UPDATE campaigns SET status = 1, updated_at = $1
		WHERE start_date::timestamp::date = now()::date`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		logrus.Debug("Update Status Base on Start Date: ", err)

		return err
	}

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
	query := `SELECT id, name, description, start_date, end_date, status, updated_at, created_at,
		metadata FROM campaigns WHERE id IS NOT NULL`

	if payload["page"].(int) > 0 || payload["limit"].(int) > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d",
			payload["limit"].(int), ((payload["page"].(int) - 1) * payload["limit"].(int)))
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

	query += where + " ORDER BY status DESC, end_date ASC" + paging
	res, err := m.getCampaign(c, query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return res, err

}

func (m *psqlCampaignRepository) getCampaign(c echo.Context, query string) ([]*models.Campaign, error) {
	var campaigns []*models.Campaign
	var campaign models.Campaign
	var rewards []models.Reward

	rows, err := m.dbBun.QueryContext(c.Request().Context(), query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		campaign = models.Campaign{}
		if err := m.dbBun.ScanRow(c.Request().Context(), rows, &campaign); err != nil {
			logger.Make(c, nil).Debug(err)

			return nil, err
		}

		// get rewards
		rewards, err = m.rwdRepo.GetRewardByCampaign(c, campaign.ID)

		if err != nil {
			logger.Make(c, nil).Debug(err)

			return nil, err
		}

		campaign.Rewards = &rewards
		campaigns = append(campaigns, &campaign)

	}

	return campaigns, nil
}

func (m *psqlCampaignRepository) GetReferralCampaign(c echo.Context, pv models.PayloadValidator) *[]*models.Campaign {
	query := `SELECT distinct c.id, c.name, c.description, c.start_date, c.end_date, c.status,
		c.updated_at, c.created_at
		FROM campaigns c
		LEFT JOIN rewards r ON c.id = r.campaign_id
		left join referral_codes rc on c.id = rc.campaign_id
		WHERE c.status = ?0 AND c.metadata->>'isReferral' = 'true' AND r.is_promo_code = ?1
		AND c.start_date::date <= ?2 AND (c.end_date::date >= ?2 OR c.end_date IS null)
		AND lower(rc.referral_code) = ?3
		order by c.created_at desc limit 1`

	rows, err := m.dbBun.QueryContext(c.Request().Context(), query, models.CampaignActive,
		models.IsPromoCodeFalse, pv.TransactionDate, strings.ToLower(pv.PromoCode))

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return nil
	}

	var campaigns []*models.Campaign
	err = m.dbBun.ScanRows(c.Request().Context(), rows, &campaigns)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return nil
	}

	if campaigns == nil {
		return nil
	}

	rewards, err := m.rwdRepo.GetRewardByCampaign(c, campaigns[0].ID)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return nil
	}

	campaigns[0].Rewards = &rewards

	return &campaigns
}

func (m *psqlCampaignRepository) GetCampaignAvailable(c echo.Context, pv models.PayloadValidator) ([]*models.Campaign, error) {
	promoCode := strings.ToLower(pv.PromoCode)

	query := fmt.Sprintf(`SELECT distinct c.id, c.name, c.description, c.start_date,
		c.end_date, c.status, c.updated_at, c.created_at, c.metadata
		FROM campaigns c
		LEFT JOIN rewards r ON c.id = r.campaign_id
		LEFT JOIN reward_tags rt ON r.id = rt.reward_id
		LEFT JOIN tags t ON rt.tag_id = t.id
		WHERE c.metadata is null and c.status = 1
		AND (
			(r.is_promo_code = 1 and (LOWER(r.promo_code) = '%s' OR lower(t.name) = '%s'))
			or (r.is_promo_code = 0)
		)
		AND c.start_date::date <= '%s'
		AND (c.end_date::date >= '%s' OR c.end_date IS null)`,
		promoCode, promoCode, pv.TransactionDate, pv.TransactionDate)

	res, err := m.getCampaign(c, query)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return nil, err
	}

	return res, err
}

func (m *psqlCampaignRepository) SavePoint(c echo.Context, cmpgnTrx *models.CampaignTrx) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	var id int64
	var query string

	if cmpgnTrx.TransactionType == models.TransactionPointTypeDebet {
		query = `INSERT INTO campaign_transactions (user_id, point_amount, transaction_type, transaction_date, ref_core, campaign_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING id`
		id = cmpgnTrx.Campaign.ID
	}

	if cmpgnTrx.TransactionType == models.TransactionPointTypeKredit {
		query = `INSERT INTO campaign_transactions (user_id, point_amount, transaction_type, transaction_date, ref_core, voucher_code_id, created_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING id`
		id = cmpgnTrx.VoucherCode.ID
	}

	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	cmpgnTrx.CreatedAt = &now
	var lastID int64
	err = stmt.QueryRow(cmpgnTrx.CIF, *cmpgnTrx.PointAmount, cmpgnTrx.TransactionType, cmpgnTrx.TransactionDate, cmpgnTrx.RefCore, id, cmpgnTrx.CreatedAt).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	cmpgnTrx.ID = lastID
	return nil
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
	var createDate, updateDate pq.NullTime
	var metadata json.RawMessage
	result := new(models.Campaign)

	query := `SELECT id, name, description, start_date, end_date, status, updated_at, created_at, metadata FROM campaigns WHERE id = $1`

	err := m.Conn.QueryRow(query, id).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.StartDate,
		&result.EndDate,
		&result.Status,
		&createDate,
		&updateDate,
		&metadata,
	)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	err = json.Unmarshal([]byte(metadata), &result.Metadata)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	result.CreatedAt = &createDate.Time
	result.UpdatedAt = &updateDate.Time

	return result, nil
}

func (m *psqlCampaignRepository) Delete(c echo.Context, id int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `DELETE FROM campaigns WHERE ID = $1`
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

	defer result.Close()
	requestLogger.Debug("Result delete campaign: ", result)

	return nil
}
