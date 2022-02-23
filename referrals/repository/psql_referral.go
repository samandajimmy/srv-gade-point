package repository

import (
	"database/sql"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"time"

	"github.com/labstack/echo"
)

type psqlReferralsRepository struct {
	Conn *sql.DB
}

// NewPsqlReferralRepository will create an object that represent the referrals.Repository interface
func NewPsqlReferralRepository(Conn *sql.DB) referrals.Repository {
	return &psqlReferralsRepository{Conn}
}

func (m *psqlReferralsRepository) CreateReferralCodes(c echo.Context, refcodes *models.ReferralCodes) error {

	var lastID int64

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	now := time.Now()
	query := `INSERT INTO referral_codes (cif, referral_code, campaign_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`
	stmt, err := m.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(refcodes.CIF, refcodes.ReferralCode, refcodes.CampaignId, &now, &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	refcodes.CreatedAt = &now
	refcodes.UpdatedAt = &now

	return nil
}

func (m *psqlReferralsRepository) GetReferralCodesByCif(c echo.Context, refCodes *models.ReferralCodes) error {

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var CreatedAt, UpdatedAt time.Time

	query := `select cif, referral_code, campaign_id, created_at, updated_at from referral_codes rc where cif = $1 order by created_at desc limit 1;`

	stmt, err := m.Conn.Query(query, refCodes.CIF)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	defer stmt.Close()

	for stmt.Next() {
		err = stmt.Scan(
			&refCodes.CIF,
			&refCodes.ReferralCode,
			&refCodes.CampaignId,
			&CreatedAt,
			&UpdatedAt,
		)

		if err != nil {
			requestLogger.Debug(err)

			return err
		}
	}

	refCodes.CreatedAt = &CreatedAt
	refCodes.UpdatedAt = &UpdatedAt

	return nil
}

func (m *psqlReferralsRepository) GetCampaignByPrefix(c echo.Context, prefix string) (int64, error) {

	var code int64
	now := time.Now()

	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	stmt, err := m.Conn.Query(`select id from campaigns c where metadata->>'prefix' = $1 and start_date <= $2 and end_date >= $3;`, prefix, &now, &now)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	defer stmt.Close()

	if !stmt.Next() {
		requestLogger.Debug(models.ErrCampaignNotReferral)

		return 0, models.ErrCampaignNotReferral
	}

	err = stmt.Scan(&code)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	return code, nil
}
