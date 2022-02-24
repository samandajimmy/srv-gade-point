package repository

import (
	"database/sql"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"
	"time"

	"github.com/labstack/echo"
	"github.com/uptrace/bun"
)

type psqlReferralsRepository struct {
	Conn *sql.DB
	Bun  *bun.DB
}

// NewPsqlReferralRepository will create an object that represent the referrals.Repository interface
func NewPsqlReferralRepository(Conn *sql.DB, Bun *bun.DB) referrals.Repository {
	return &psqlReferralsRepository{Conn, Bun}
}

func (m *psqlReferralsRepository) CreateReferral(c echo.Context, refcodes models.ReferralCodes) (models.ReferralCodes, error) {

	now := time.Now()
	refcodes.CreatedAt = now
	refcodes.UpdatedAt = now

	query := `INSERT INTO referral_codes (cif, referral_code, campaign_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?) RETURNING id`

	_, err := m.Bun.QueryContext(c.Request().Context(), query, refcodes.CIF, refcodes.ReferralCode, refcodes.CampaignId, refcodes.CreatedAt, refcodes.UpdatedAt)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ReferralCodes{}, err
	}

	return refcodes, nil
}

func (m *psqlReferralsRepository) GetReferralByCif(c echo.Context, refCodes models.ReferralCodes) (models.ReferralCodes, error) {

	var result models.ReferralCodes

	query := `select cif, referral_code, campaign_id, created_at, updated_at from referral_codes rc where cif = ? order by created_at desc limit 1;`

	rows, err := m.Bun.QueryContext(c.Request().Context(), query, refCodes.CIF)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ReferralCodes{}, err
	}

	err = m.Bun.ScanRows(c.Request().Context(), rows, &result)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return models.ReferralCodes{}, err
	}

	return result, nil
}

func (m *psqlReferralsRepository) GetCampaignId(c echo.Context, prefix string) (int64, error) {

	var code int64
	now := time.Now()

	query := `select id from campaigns c where metadata->>'prefix' = ? and start_date <= ? and end_date >= ?;`

	rows, err := m.Bun.QueryContext(c.Request().Context(), query, prefix, &now, &now)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	err = m.Bun.ScanRows(c.Request().Context(), rows, &code)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return 0, err
	}

	return code, nil
}
