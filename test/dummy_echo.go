package test

import (
	"gade/srv-gade-point/config"
	"gade/srv-gade-point/database"
	"gade/srv-gade-point/helper"
	"gade/srv-gade-point/middleware"
	"gade/srv-gade-point/mocks"
	"gade/srv-gade-point/models"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang/mock/gomock"
	"github.com/labstack/echo"
)

type (
	DummyEcho struct {
		EchoObj  *echo.Echo
		Request  *http.Request
		Response *httptest.ResponseRecorder
		Context  echo.Context
	}

	customValidator struct {
		middleware.CustomValidator
	}
)

func NewDummyEcho(method, path string, pl ...interface{}) DummyEcho {
	var body string
	e := echo.New()
	cv := customValidator{}
	cv.CustomValidation()
	e.Validator = &cv

	if pl != nil {
		body = helper.ToJson(pl[0])
	}

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)

	return DummyEcho{
		EchoObj:  e,
		Request:  req,
		Response: resp,
		Context:  c,
	}
}

func NewDbConn() (*database.DbConfig, *migrate.Migrate) {
	dbConfig := database.DbConfig{
		Host:     os.Getenv(`DB_TEST_HOST`),
		Port:     os.Getenv(`DB_TEST_PORT`),
		Username: os.Getenv(`DB_TEST_USER`),
		Password: os.Getenv(`DB_TEST_PASS`),
		Name:     os.Getenv(`DB_TEST_NAME`),
		Logger:   os.Getenv(`DB_TEST_LOGGER`),
	}

	db := database.NewDbConn(dbConfig)
	migrator := db.Migrate(db.Sql)

	return db, migrator
}

func LoadRealRepoUsecase(db *database.DbConfig) (config.Repositories, config.Usecases) {
	repos := config.NewRepositories(db.Sql, db.Bun)
	usecase := config.NewUsecases(repos)

	return repos, usecase
}

func (de *DummyEcho) LoadMockRepoUsecase(mockCtrl *gomock.Controller) (mocks.MockRepositories, mocks.MockUsecases) {
	mockRepos := mocks.NewMockRepository(mockCtrl)
	mockUsecase := mocks.NewMockUsecases(mockCtrl)

	return mockRepos, mockUsecase
}

func (cv *customValidator) Validate(i interface{}) error {
	rValue := reflect.ValueOf(i)
	if rValue.Kind() == reflect.Ptr {
		rValue = rValue.Elem()
	}

	if rValue.Kind() == reflect.Struct {
		return cv.Validator.Struct(i)
	}

	obj, ok := i.(map[string]interface{})

	if !ok {
		return models.ErrInternalServerError
	}

	isError, ok := obj["isError"].(bool)

	if !ok {
		return nil
	}

	if isError {
		return models.ErrInternalServerError
	}

	return nil
}
