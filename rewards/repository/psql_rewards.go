package repository

import (
	"database/sql"
	"encoding/json"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/rewards"
	"time"

	"github.com/labstack/echo"
)

type psqlRewardRepository struct {
	Conn *sql.DB
}

// NewPsqlRewardRepository will create an object that represent the rewards.Repository interface
func NewPsqlRewardRepository(Conn *sql.DB) rewards.Repository {
	return &psqlRewardRepository{Conn}
}

func (rwdRepo *psqlRewardRepository) CreateReward(c echo.Context, reward *models.Reward, campaignID int64) error {
	var lastID int64
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `INSERT INTO rewards (name, description, terms_and_conditions, how_to_use,
		journal_account, promo_code, is_promo_code, custom_period, type, validators,
		campaign_id, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11,
		$12) RETURNING id`
	stmt, err := rwdRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	validator, err := json.Marshal(reward.Validators)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(
		reward.Name, reward.Description, reward.TermsAndConditions, reward.HowToUse,
		reward.JournalAccount, reward.PromoCode, reward.IsPromoCode, reward.CustomPeriod,
		reward.Type, string(validator), campaignID, &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	reward.ID = lastID
	reward.CreatedAt = &now
	return nil
}

func (rwdRepo *psqlRewardRepository) DeleteByCampaign(c echo.Context, campaignID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `DELETE FROM rewards WHERE campaign_id = $1`
	stmt, err := rwdRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	result, err := stmt.Query(campaignID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	requestLogger.Debug("rewards deleted: ", result)

	return nil
}

func (rwdRepo *psqlRewardRepository) CreateRewardTag(c echo.Context, tag *models.Tag, rewardID int64) error {
	var lastID int64
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `INSERT INTO reward_tags (reward_id, tag_id, created_at) VALUES ($1, $2, $3) RETURNING id`
	stmt, err := rwdRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(tag.ID, rewardID, &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	return nil
}

func (rwdRepo *psqlRewardRepository) DeleteRewardTag(c echo.Context, rewardID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `DELETE FROM reward_tags WHERE reward_id = $1`
	stmt, err := rwdRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	result, err := stmt.Query(rewardID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	requestLogger.Debug("reward_tags deleted: ", result)

	return nil
}
