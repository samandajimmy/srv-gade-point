package users

import (
	"context"
	"gade/srv-gade-point/models"
)

// UseCase represent the user's usecases
type UseCase interface {
	CreateUser(ctx context.Context, usr *models.User) error
	ValidateUser(ctx context.Context, usr models.User) (*models.User, error)
}
