package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"gade/srv-gade-point/database"
	gcdb "gade/srv-gade-point/database"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/quotas"
	"gade/srv-gade-point/rewards"
	"gade/srv-gade-point/tags"
	"math/rand"
	"time"

	"github.com/labstack/echo"
	"github.com/lib/pq"
)

type psqlRewardRepository struct {
	Conn     *sql.DB
	dbBun    *database.DbBun
	quotRepo quotas.Repository
	tagRepo  tags.Repository
}

// NewPsqlRewardRepository will create an object that represent the rewards.Repository interface
func NewPsqlRewardRepository(Conn *sql.DB, dbBun *database.DbBun, quotRepo quotas.Repository, tagRepo tags.Repository) rewards.RRepository {
	return &psqlRewardRepository{Conn, dbBun, quotRepo, tagRepo}
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

	defer result.Close()
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

	defer result.Close()
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

	if err != nil {
		requestLogger.Debug(err)

		return rewards, err
	}

	defer rows.Close()

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

	if err != nil {
		requestLogger.Debug(err)

		return reward, err
	}

	defer rows.Close()

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

func (rwdRepo *psqlRewardRepository) CountRewards(c echo.Context, rewardPayload *models.RewardsPayload) (int64, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var count int64

	query := `SELECT count(ID) FROM rewards`

	err := rwdRepo.Conn.QueryRow(query).Scan(&count)

	if err != nil {
		requestLogger.Debug(err)

		return 0, err
	}

	return count, nil
}

func (rwdRepo *psqlRewardRepository) GetRewards(c echo.Context, rewardPayload *models.RewardsPayload) ([]models.Reward, error) {
	var data []models.Reward
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	var paging, where string

	query := `SELECT r.id, r.name, r.promo_code, r.campaign_id, c.name FROM rewards r join campaigns c on r.campaign_id = c.id`

	if rewardPayload.Page > 0 || rewardPayload.Limit > 0 {
		paging = fmt.Sprintf(" LIMIT %d OFFSET %d", rewardPayload.Limit, ((rewardPayload.Page - 1) * rewardPayload.Limit))
	}

	query += where + " order by r.updated_at desc" + paging + ";"
	rows, err := rwdRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var rwd models.Reward
		rwd.Campaign = &models.Campaign{}

		err = rows.Scan(
			&rwd.ID,
			&rwd.Name,
			&rwd.PromoCode,
			&rwd.CampaignID,
			&rwd.Campaign.Name,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		data = append(data, rwd)
	}

	return data, nil
}

func (rwdRepo *psqlRewardRepository) RGetRandomId(n int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, n)

	for i := range b {
		b[i] = models.LetterBytes[rand.Int63()%int64(len(models.LetterBytes))]
	}

	return string(b)
}

func (rwdRepo *psqlRewardRepository) GetRewardPromotions(c echo.Context, pv models.RewardPromotionLists) ([]*models.RewardPromotions, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	where := ""

	query := `SELECT r.id, r.name, r.description, r.terms_and_conditions, r.how_to_use, 
		r.promo_code, r.validators->>'product', r.validators->>'transactionType', 
		r.validators->>'minTransactionAmount', COALESCE(r.validators->>'isPrivate','0')
		FROM rewards r
		LEFT JOIN campaigns c on r.campaign_id = c.id
		WHERE c.status = 1
  		AND r.is_promo_code = 1
  		AND c.start_date::date <= now()::date
		AND (c.end_date::date >= now()::date OR c.end_date IS null)`

	if pv.Product != "" {
		where += " AND r.validators->>'product' LIKE '%" + pv.Product + "%'"
	}

	if pv.TransactionType != "" {
		where += " AND r.validators->>'transactionType' LIKE '%" + pv.TransactionType + "%'"
	}

	if pv.Channel != "" {
		where += " AND r.validators->>'channel' LIKE '%" + pv.Channel + "%'"
	}

	query += where
	res, err := rwdRepo.getRewardPromotions(c, query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	return res, err
}

func (rwdRepo *psqlRewardRepository) getRewardPromotions(c echo.Context, query string) ([]*models.RewardPromotions, error) {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	result := make([]*models.RewardPromotions, 0)

	rows, err := rwdRepo.Conn.Query(query)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		t := new(models.RewardPromotions)
		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.Description,
			&t.TermsAndConditions,
			&t.HowToUse,
			&t.PromoCode,
			&t.Product,
			&t.TransactionType,
			&t.MinTransactionAmount,
			&t.IsPrivate,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		result = append(result, t)
	}

	return result, nil
}

func (rwdRepo *psqlRewardRepository) RGetRewardDetail(c echo.Context, rewardId int64) (models.Reward, error) {
	var result models.Reward

	query := `select  id, name, description, terms_and_conditions, how_to_use, journal_account, 
			promo_code, is_promo_code, custom_period, type, validators, campaign_id, 
			updated_at, created_at from rewards where id = ?`

	err := rwdRepo.dbBun.QueryThenScan(c, &result, query, rewardId)

	if err != nil {
		logger.Make(c, nil).Debug(err)

		return result, err
	}

	return result, nil
}

func (rwdRepo *psqlRewardRepository) RPostCoreTrx(c echo.Context, coreTrx []models.CoreTrxPayload) error {
	var nilFilters []string
	createdAt := time.Now()
	stmts := []*gcdb.PipelineStmt{}
	for _, trx := range coreTrx {

		stmts = append(stmts, gcdb.NewPipelineStmt(`INSERT INTO core_transactions 
		(created_at, transaction_amount, loan_amount, interest_amount, product_code, 
		transaction_date, total_reward, transaction_id, marketing_code, transaction_type,
		inq_status, root_ref_trx) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
			nilFilters, createdAt, trx.TrxAmount, trx.LoanAmount,
			trx.InterestAmount, trx.ProductCode, trx.TrxDate, trx.RwdTotal, trx.TrxID,
			trx.MarketingCode, trx.TrxType, trx.InqStatus, trx.RootRefTrx))
	}

	err := gcdb.WithTransaction(rwdRepo.Conn, func(tx gcdb.Transaction) error {
		return gcdb.RunPipelineQueryRow(tx, stmts...)
	})

	if err != nil {
		logger.Make(c, nil).Debug(err)
		return err
	}

	return nil
}
