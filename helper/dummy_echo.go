package helper

import (
	"gade/srv-gade-point/models"
	"net/http"
	"net/http/httptest"
	"strings"

	"gopkg.in/go-playground/validator.v9"

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
		validator *validator.Validate
	}
)

func NewDummyEcho(method, path string, pl ...interface{}) DummyEcho {
	var body string
	e := echo.New()
	e.Validator = &customValidator{validator: validator.New()}

	if pl != nil {
		body = ToJson(pl[0])
	}

	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)

	return DummyEcho{e, req, resp, c}
}

func (cv *customValidator) Validate(i interface{}) error {

	_, ok := i.(map[string]interface{})

	if ok {
		return nil
	}

	return models.ErrInternalServerError
}
