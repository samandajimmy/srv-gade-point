package users

import (
	"context"
	"gade/srv-gade-point/models"
)

// Repository represent the user's repository contract
type Repository interface {
	Create(ctx context.Context, usr *models.User) error
	GetByUsername(ctx context.Context, usr *models.User) error
}
