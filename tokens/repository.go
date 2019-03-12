package tokens

import (
	"context"
	"gade/srv-gade-point/models"
)

// Repository represent the Account Token's repository contract
type Repository interface {
	Create(ctx context.Context, accToken *models.AccountToken) error
	GetByUsername(ctx context.Context, accToken *models.AccountToken) error
	UpdateToken(ctx context.Context, accToken *models.AccountToken) error
}
