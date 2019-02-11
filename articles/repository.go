package articles

import (
	"context"

	"gade/srv-gade-point/models"
)

// Repository represent the article's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Article, nextCursor string, err error)
	Store(ctx context.Context, a *models.Article) error
}
