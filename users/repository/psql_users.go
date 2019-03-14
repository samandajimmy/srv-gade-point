package repository

import (
	"context"
	"database/sql"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/users"
	"time"

	"github.com/lib/pq"
)

type psqlUserRepository struct {
	Conn *sql.DB
}

// NewPsqlUserRepository will create an object that represent the user.Repository interface
func NewPsqlUserRepository(Conn *sql.DB) users.Repository {
	return &psqlUserRepository{Conn}
}

func (m *psqlUserRepository) Create(ctx context.Context, usr *models.User) error {
	var lastID int64
	now := time.Now()
	defStatus := int8(1)
	defRole := int8(0)

	usr.CreatedAt = &now
	usr.Status = &defStatus
	usr.Role = &defRole

	query := `INSERT INTO users (username, email, password, status, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)  RETURNING id`
	stmt, err := m.Conn.PrepareContext(ctx, query)

	if err != nil {
		return err
	}

	err = stmt.QueryRowContext(ctx, usr.Username, usr.Email, usr.Password,
		usr.Status, usr.Role, usr.CreatedAt).Scan(&lastID)

	if err != nil {
		return err
	}

	usr.ID = lastID
	return nil
}

func (m *psqlUserRepository) GetByUsername(ctx context.Context, usr *models.User) error {
	var updatedAt, createdAt pq.NullTime
	query := `SELECT id, username, email, password, status, role, updated_at, created_at
		FROM users
		WHERE status = 1 AND username = $1`

	err := m.Conn.QueryRowContext(ctx, query, usr.Username).Scan(
		&usr.ID, &usr.Username, &usr.Email, &usr.Password,
		&usr.Status, &usr.Role, &updatedAt, &createdAt,
	)

	if err != nil {
		return err
	}

	usr.CreatedAt = &createdAt.Time
	usr.UpdatedAt = &updatedAt.Time

	return nil
}
