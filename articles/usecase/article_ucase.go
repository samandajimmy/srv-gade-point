package usecase

import (
	"context"
	"time"

	"gade/srv-gade-point/articles"
	"gade/srv-gade-point/models"
)

type articleUsecase struct {
	articleRepo    articles.Repository
	contextTimeout time.Duration
}

// NewArticleUsecase will create new an articleUsecase object representation of articles.Usecase interface
func NewArticleUsecase(a articles.Repository, timeout time.Duration) articles.Usecase {
	return &articleUsecase{
		articleRepo:    a,
		contextTimeout: timeout,
	}
}

/*
* In this function below, I'm using errgroup with the pipeline pattern
* Look how this works in this package explanation
* in godoc: https://godoc.org/golang.org/x/sync/errgroup#ex-Group--Pipeline
 */

func (a *articleUsecase) Fetch(c context.Context, cursor string, num int64) ([]*models.Article, string, error) {
	if num == 0 {
		num = 10
	}

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	listArticle, nextCursor, err := a.articleRepo.Fetch(ctx, cursor, num)
	if err != nil {
		return nil, "", err
	}

	if err != nil {
		return nil, "", err
	}

	return listArticle, nextCursor, nil
}

func (a *articleUsecase) Store(c context.Context, m *models.Article) error {

	ctx, cancel := context.WithTimeout(c, a.contextTimeout)
	defer cancel()

	err := a.articleRepo.Store(ctx, m)
	if err != nil {
		return err
	}
	return nil
}
