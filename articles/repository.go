package articles

import (
	"context"

	"github.com/gade-dev/gade-srv-boilerplate-go/models"
)

// Repository represent the article's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Article, nextCursor string, err error)
	Store(ctx context.Context, a *models.Article) error
}
