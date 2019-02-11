package http

import (
	"context"
	"net/http"
	"strconv"

	"gade/srv-gade-point/articles"
	"gade/srv-gade-point/models"
	"github.com/labstack/echo"
	"github.com/sirupsen/logrus"
	validator "gopkg.in/go-playground/validator.v9"
)

// ResponseError represent the response error struct
type ResponseError struct {
	Message string `json:"message"`
}

// ArticlesHandler represent the httphandler for articles
type ArticlesHandler struct {
	AUsecase articles.Usecase
}

// NewArticlesHandler represent to register articles endpoint
func NewArticlesHandler(e *echo.Echo, us articles.Usecase) {
	handler := &ArticlesHandler{
		AUsecase: us,
	}

	e.GET("/articles", handler.FetchArticles)
	e.POST("/articles", handler.Store)
}

func (a *ArticlesHandler) FetchArticles(c echo.Context) error {

	numS := c.QueryParam("num")
	num, _ := strconv.Atoi(numS)
	cursor := c.QueryParam("cursor")
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}
	listAr, nextCursor, err := a.AUsecase.Fetch(ctx, cursor, int64(num))

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	c.Response().Header().Set(`X-Cursor`, nextCursor)
	return c.JSON(http.StatusOK, listAr)
}

func (a *ArticlesHandler) Store(c echo.Context) error {
	var article models.Article
	err := c.Bind(&article)
	if err != nil {
		return c.JSON(http.StatusUnprocessableEntity, err.Error())
	}

	if ok, err := isRequestValid(&article); !ok {
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	err = a.AUsecase.Store(ctx, &article)

	if err != nil {
		return c.JSON(getStatusCode(err), ResponseError{Message: err.Error()})
	}
	return c.JSON(http.StatusCreated, article)
}

func isRequestValid(m *models.Article) (bool, error) {

	validate := validator.New()

	err := validate.Struct(m)
	if err != nil {
		return false, err
	}
	return true, nil
}

func getStatusCode(err error) int {

	if err == nil {
		return http.StatusOK
	}
	logrus.Error(err)
	switch err {
	case models.ErrInternalServerError:
		return http.StatusInternalServerError
	case models.ErrNotFound:
		return http.StatusNotFound
	case models.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}
