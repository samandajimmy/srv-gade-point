package middleware

import (
	"gade/srv-gade-point/models"
	"os"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type customMiddleware struct {
	e *echo.Echo
}

var echGroup models.EchoGroup

// InitMiddleware to generate all middleware that app need
func InitMiddleware(ech *echo.Echo, echoGroup models.EchoGroup) {
	cm := &customMiddleware{ech}
	echGroup = echoGroup
	ech.Use(middleware.RequestIDWithConfig(middleware.DefaultRequestIDConfig))

	ech.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "requestID=${id}, method=${method}, status=${status}, uri=${uri}, latency=${latency_human} " +
			"host=${host}, remote_ip=${remote_ip}, user_agent=${user_agent}, error=${error} \n",
	}))

	ech.Use(middleware.Recover())
	cm.cors()
	cm.basicAuth()
	cm.jwtAuth()
}

func (cm customMiddleware) cors() {
	cm.e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"Access-Control-Allow-Origin"},
		AllowMethods: []string{"*"},
	}))
}

func (cm customMiddleware) basicAuth() {
	echGroup.Token.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == os.Getenv(`BASIC_USERNAME`) && password == os.Getenv(`BASIC_PASSWORD`) {
			return true, nil
		}

		return false, nil
	}))
}

func (cm customMiddleware) jwtAuth() {
	echGroup.Admin.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(os.Getenv(`JWT_SECRET`)),
	}))

	echGroup.API.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(os.Getenv(`JWT_SECRET`)),
	}))

}
