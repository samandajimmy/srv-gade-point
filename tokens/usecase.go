package tokens

import (
	"context"
	"gade/srv-gade-point/models"
)

// UseCase represent the token's usecases
type UseCase interface {
	CreateToken(ctx context.Context, accToken *models.AccountToken) error
	GetToken(ctx context.Context, username string, password string) (*models.AccountToken, error)
	RefreshToken(ctx context.Context, username string, password string) (*models.AccountToken, error)
}
