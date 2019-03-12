package repository

import (
	"context"
	"database/sql"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/tokens"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/gommon/log"
	"github.com/lib/pq"
)

type psqlTokenRepository struct {
	Conn *sql.DB
}

// NewPsqlTokenRepository will create an object that represent the tokens.Repository interface
func NewPsqlTokenRepository(Conn *sql.DB) tokens.Repository {
	return &psqlTokenRepository{Conn}
}

func (m *psqlTokenRepository) Create(ctx context.Context, accToken *models.AccountToken) error {
	var lastID int64
	now := time.Now()
	tokenExp := now.Add(stringToDuration(os.Getenv(`JWT_TOKEN_EXP`)) * time.Hour)
	defStatus := int8(1)

	token, err := createJWTToken(accToken, now, tokenExp)

	if err != nil {
		log.Error(err)

		return err
	}

	accToken.Token = token
	accToken.ExpireAt = &tokenExp
	accToken.CreatedAt = &now
	accToken.Status = &defStatus

	query := `INSERT INTO account_tokens (username, password, token, expire_at, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)  RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	err = stmt.QueryRowContext(ctx, accToken.Username, accToken.Password, accToken.Token,
		accToken.ExpireAt, accToken.Status, accToken.CreatedAt).Scan(&lastID)

	if err != nil {
		return err
	}

	accToken.ID = lastID
	return nil
}

func (m *psqlTokenRepository) GetByUsername(ctx context.Context, accToken *models.AccountToken) error {
	var expireAt, updatedAt, createdAt pq.NullTime
	query := `SELECT id, username, password, token, expire_at, status, updated_at, created_at
		FROM account_tokens
		WHERE status = 1 AND username = $1`

	err := m.Conn.QueryRowContext(ctx, query, accToken.Username).Scan(
		&accToken.ID, &accToken.Username, &accToken.Password,
		&accToken.Token, &expireAt, &accToken.Status,
		&updatedAt, &createdAt,
	)

	if err != nil {
		return err
	}

	accToken.ExpireAt = &expireAt.Time
	accToken.CreatedAt = &createdAt.Time
	accToken.UpdatedAt = &updatedAt.Time

	return nil
}

func (m *psqlTokenRepository) UpdateToken(ctx context.Context, accToken *models.AccountToken) error {
	var ID int64
	now := time.Now()
	tokenExp := now.Add(stringToDuration(os.Getenv(`JWT_TOKEN_EXP`)) * time.Hour)
	token, err := createJWTToken(accToken, now, tokenExp)

	if err != nil {
		log.Error(err)

		return err
	}

	query := `UPDATE account_tokens SET token = $1, expire_at = $2, updated_at = $3 WHERE username = $4 RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		log.Error(err)

		return err
	}

	err = stmt.QueryRowContext(ctx, token, tokenExp, now, accToken.Username).Scan(&ID)

	if err != nil {
		log.Error(err)

		return err
	}

	return nil
}

func stringToDuration(str string) time.Duration {
	hours, _ := strconv.Atoi(str)

	return time.Duration(hours)
}

func createJWTToken(accountToken *models.AccountToken, now time.Time, tokenExp time.Time) (string, error) {
	claims := models.Token{
		accountToken.Username,
		jwt.StandardClaims{
			Id:        accountToken.Username,
			ExpiresAt: tokenExp.Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return rawToken.SignedString([]byte(os.Getenv(`JWT_SECRET`)))
}
