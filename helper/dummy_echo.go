package helper

import (
	"net/http"
	"net/http/httptest"

	"github.com/labstack/echo"
)

type DummyEcho struct {
	EchoObj  *echo.Echo
	Request  *http.Request
	Response *httptest.ResponseRecorder
	Context  echo.Context
}

type DummyTestData struct {
	Campaign map[string]interface{}
}

func NewDummyEcho(method, path string) DummyEcho {
	e := echo.New()
	req := httptest.NewRequest(method, path, nil)
	resp := httptest.NewRecorder()
	c := e.NewContext(req, resp)

	return DummyEcho{e, req, resp, c}
}
