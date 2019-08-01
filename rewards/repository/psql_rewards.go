package repository

import (
	"database/sql"
	"encoding/json"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/quotas"
	"gade/srv-gade-point/rewards"
	"gade/srv-gade-point/tags"
	"time"

	"github.com/labstack/echo"
	"github.com/lib/pq"
)

type psqlRewardRepository struct {
	Conn     *sql.DB
	quotRepo quotas.Repository
	tagRepo  tags.Repository
}

// NewPsqlRewardRepository will create an object that represent the rewards.Repository interface
func NewPsqlRewardRepository(Conn *sql.DB, quotRepo quotas.Repository, tagRepo tags.Repository) rewards.Repository {
	return &psqlRewardRepository{Conn, quotRepo, tagRepo}
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
	query := `INSERT INTO reward_tags (reward_id, tag_id, created_at) VALUES ($1, $2, $3) RETURNING reward_id`
	stmt, err := rwdRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(rewardID, tag.ID, &now).Scan(&lastID)

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

func (rwdRepo *psqlRewardRepository) GetRewardByCampaign(c echo.Context, campaignID int64) ([]models.Reward, error) {
	var rewards []models.Reward
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT id, name, description, terms_and_conditions, how_to_use, journal_account, promo_code, is_promo_code, custom_period, type, validators,
		campaign_id, created_at, updated_at FROM rewards WHERE campaign_id = $1`
	rows, err := rwdRepo.Conn.Query(query, campaignID)
	defer rows.Close()

	if err != nil {
		requestLogger.Debug(err)

		return rewards, err
	}

	for rows.Next() {
		var reward models.Reward
		var createDate, updateDate pq.NullTime
		var validator json.RawMessage

		err = rows.Scan(
			&reward.ID, &reward.Name, &reward.Description, &reward.TermsAndConditions, &reward.HowToUse, &reward.JournalAccount, &reward.PromoCode,
			&reward.IsPromoCode, &reward.CustomPeriod, &reward.Type, &validator, &reward.CampaignID, &createDate, &updateDate,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		reward.CreatedAt = &createDate.Time
		reward.UpdatedAt = &updateDate.Time

		// get quotas
		quotas, err := rwdRepo.quotRepo.GetQuotaByReward(c, reward.ID)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		// get tags
		tags, err := rwdRepo.tagRepo.GetTagByReward(c, reward.ID)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		reward.Quotas = &quotas
		reward.Tags = &tags

		if err = json.Unmarshal([]byte(validator), &reward.Validators); err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		rewards = append(rewards, reward)
	}

	return rewards, nil
}

func (rwdRepo *psqlRewardRepository) GetRewardTags(c echo.Context, reward *models.Reward) (*models.Reward, error) {
	var tags []models.Tag
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `SELECT t.name FROM reward_tags rt JOIN tags t ON rt.tag_id = t.id  WHERE reward_id = $1`
	rows, err := rwdRepo.Conn.Query(query, reward.ID)
	defer rows.Close()

	if err != nil {
		requestLogger.Debug(err)

		return reward, err
	}

	for rows.Next() {
		var tag models.Tag

		err = rows.Scan(&tag.Name)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		tags = append(tags, tag)
	}

	reward.Tags = &tags

	return reward, nil
}
