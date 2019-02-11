package repository

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/gade-dev/srv-gade-point/articles"
	"github.com/gade-dev/srv-gade-point/models"
)

const (
	timeFormat = "2006-01-02T15:04:05.999Z07:00" // reduce precision from RFC3339Nano as date format
)

type psqlArticleRepository struct {
	Conn *sql.DB
}

// NewPsqlArticleRepository will create an object that represent the articles.Repository interface
func NewPsqlArticleRepository(Conn *sql.DB) articles.Repository {
	return &psqlArticleRepository{Conn}
}

func (m *psqlArticleRepository) Fetch(ctx context.Context, cursor string, num int64) ([]*models.Article, string, error) {
	query := `SELECT * FROM articles`

	decodedCursor, err := DecodeCursor(cursor)
	if err != nil && cursor != "" {
		return nil, "", models.ErrBadParamInput
	}
	res, err := m.fetch(ctx, query, decodedCursor, num)
	if err != nil {
		return nil, "", err
	}
	nextCursor := ""
	if len(res) == int(num) {
		nextCursor = EncodeCursor(res[len(res)-1].CreatedAt)
	}
	return res, nextCursor, err

}

func (m *psqlArticleRepository) Store(ctx context.Context, a *models.Article) error {
	// query := `INSERT  articles SET title=? , content=? , author_id=?, updated_at=? , created_at=?`

	query := `INSERT INTO articles (title, content, author_id, updated_at, created_at)
		VALUES ($1, $2, $3, $4, $5)  RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)
	if err != nil {
		return err
	}

	logrus.Debug("Created At: ", a.CreatedAt)
	var lastID int64
	err = stmt.QueryRowContext(ctx, a.Title, a.Content, a.Author.ID, a.UpdatedAt, a.CreatedAt).Scan(&lastID)
	if err != nil {
		return err
	}

	a.ID = lastID
	return nil
}

func (m *psqlArticleRepository) fetch(ctx context.Context, query string, args ...interface{}) ([]*models.Article, error) {
	rows, err := m.Conn.QueryContext(ctx, query)
	fmt.Println(err)

	if err != nil {
		logrus.Error(err)
		return nil, err
	}

	defer rows.Close()
	result := make([]*models.Article, 0)
	for rows.Next() {
		t := new(models.Article)
		authorID := int64(0)
		err = rows.Scan(
			&t.ID,
			&t.Title,
			&t.Content,
			&authorID,
			&t.UpdatedAt,
			&t.CreatedAt,
		)

		if err != nil {
			logrus.Error(err)
			return nil, err
		}
		t.Author = models.Author{
			ID: authorID,
		}
		result = append(result, t)
	}

	return result, nil
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
