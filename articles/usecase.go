package articles

import (
	"context"

	"github.com/gade-dev/srv-gade-point/models"
)

// Usecase represent the article's usecases
type Usecase interface {
	Fetch(ctx context.Context, cursor string, num int64) ([]*models.Article, string, error)
	Store(context.Context, *models.Article) error
}
