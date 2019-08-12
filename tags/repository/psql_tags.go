package repository

import (
	"database/sql"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/tags"
	"time"

	"github.com/labstack/echo"
)

type psqlTagRepository struct {
	Conn *sql.DB
}

// NewPsqlTagRepository will create an object that represent the tags.Repository interface
func NewPsqlTagRepository(Conn *sql.DB) tags.Repository {
	return &psqlTagRepository{Conn}
}

func (tgRepo *psqlTagRepository) CreateTag(c echo.Context, tag *models.Tag, campaignID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	now := time.Now()
	query := `INSERT INTO tags (name, created_at) VALUES ($1, $2) RETURNING id`
	stmt, err := tgRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	var lastID int64

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	err = stmt.QueryRow(
		tag.Name, &now).Scan(&lastID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	tag.ID = lastID
	tag.CreatedAt = &now
	return nil
}

func (tgRepo *psqlTagRepository) Delete(c echo.Context, tagID int64) error {
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)
	query := `DELETE FROM tags WHERE id = $1`
	stmt, err := tgRepo.Conn.Prepare(query)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	result, err := stmt.Query(tagID)

	if err != nil {
		requestLogger.Debug(err)

		return err
	}

	defer result.Close()
	requestLogger.Debug("tags deleted: ", result)

	return nil
}

func (tgRepo *psqlTagRepository) GetTagByReward(c echo.Context, rewardID int64) ([]models.Tag, error) {
	var result []models.Tag
	logger := models.RequestLogger{}
	requestLogger := logger.GetRequestLogger(c, nil)

	query := `SELECT t.id, t.name, t.created_at, t.updated_at FROM reward_tags rt join tags t on rt.tag_id = t.id where rt.reward_id = $1`
	rows, err := tgRepo.Conn.Query(query, rewardID)

	if err != nil {
		requestLogger.Debug(err)

		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var t models.Tag

		err = rows.Scan(
			&t.ID,
			&t.Name,
			&t.CreatedAt,
			&t.UpdatedAt,
		)

		if err != nil {
			requestLogger.Debug(err)

			return nil, err
		}

		result = append(result, t)
	}

	return result, err
}
