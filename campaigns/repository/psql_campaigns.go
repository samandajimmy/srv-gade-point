package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/sirupsen/logrus"

	"gade/srv-gade-point/campaigns"
	"gade/srv-gade-point/models"
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

	query := `INSERT INTO campaigns (name, description, start_date, end_date, status, validators, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)  RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	logrus.Debug("Created At: ", a.CreatedAt)

	var lastID int64

	validator, err := json.Marshal(a.Validators)
	if err != nil {
		return err
	}

	err = stmt.QueryRowContext(ctx, a.Name, a.Description, a.StartDate, a.EndDate, a.Status, string(validator), a.CreatedAt).Scan(&lastID)
	if err != nil {
		return err
	}

	a.ID = lastID
	return nil
}

func DecodeCursor(encodedTime string) (time.Time, error) {
	byt, err := base64.StdEncoding.DecodeString(encodedTime)
	if err != nil {
		return time.Time{}, err
	}

	timeString := string(byt)
	t, err := time.Parse(timeFormat, timeString)

	return t, err
}

func EncodeCursor(t time.Time) string {
	timeString := t.Format(timeFormat)

	return base64.StdEncoding.EncodeToString([]byte(timeString))
}
