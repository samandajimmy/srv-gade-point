package repository

import (
	"context"
	"database/sql"
	"encoding/json"
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

type JSONB []byte

// NewPsqlCampaignRepository will create an object that represent the campaigns.Repository interface
func NewPsqlCampaignRepository(Conn *sql.DB) campaigns.Repository {
	return &psqlCampaignRepository{Conn}
}

func (m *psqlCampaignRepository) CreateCampaign(ctx context.Context, a *models.Campaign) error {

	query := `INSERT INTO campaigns (name, description, start_date, end_date, status, validators, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING id`
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

	err = stmt.QueryRowContext(ctx, a.Name, a.Description, a.StartDate, a.EndDate, a.Status, string(validator), time.Now()).Scan(&lastID)
	if err != nil {
		return err
	}

	a.ID = lastID
	a.CreatedAt = time.Now()
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

func (m *psqlCampaignRepository) GetCampaign(ctx context.Context, name string, status string, startDate string, endDate string) ([]*models.Campaign, error) {
	query := `SELECT * FROM campaigns WHERE id IS NOT NULL`

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

	res, err := m.getCampaign(ctx, query)
	if err != nil {
		return nil, err
	}

	return res, err

}

func (m *psqlCampaignRepository) getCampaign(ctx context.Context, query string) ([]*models.Campaign, error) {
	var validator json.RawMessage
	rows, err := m.Conn.QueryContext(ctx, query)
	if err != nil {
		logrus.Error(err)
		return nil, err
	}
	defer rows.Close()

	result := make([]*models.Campaign, 0)
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
