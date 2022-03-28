package middleware

import (
	"encoding/json"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"

	"net/url"
	"os"
	"reflect"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type CustomValidator struct {
	Validator *validator.Validate
}

type customMiddleware struct {
	e *echo.Echo
}

var (
	echGroup  models.EchoGroup
	jwtMiddFn echo.MiddlewareFunc
)

// InitMiddleware to generate all middleware that app need
func InitMiddleware(ech *echo.Echo, echoGroup models.EchoGroup) {
	jwtMiddFn = middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte(os.Getenv(`JWT_SECRET`)),
	})

	cm := &customMiddleware{ech}
	echGroup = echoGroup
	ech.Use(middleware.RequestIDWithConfig(middleware.DefaultRequestIDConfig))
	cm.customBodyDump()

	ech.Use(middleware.Recover())
	cm.cors()
	cm.basicAuth()
	cm.jwtAuth()
	cv := CustomValidator{}
	cm.e.Validator = &cv
}

func ReferralAuth(handlerFn echo.HandlerFunc) {
	echGroup.Referral.Use(
		jwtMiddFn,
		func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				return handlerFn(c)
			}
		},
	)
}

func (cm *customMiddleware) customBodyDump() {
	cm.e.Use(middleware.BodyDumpWithConfig(middleware.BodyDumpConfig{
		Handler: func(c echo.Context, req, resp []byte) {
			bodyParser(c, &req)
			reqBody := c.Request()

			logger.MakeWithoutReportCaller(c, req).Info("Request payload for endpoint " + reqBody.Method + " " + reqBody.URL.Path)
			logger.MakeWithoutReportCaller(c, resp).Info("Response payload for endpoint " + reqBody.Method + " " + reqBody.URL.Path)
		},
	}))
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
	echGroup.Admin.Use(jwtMiddFn)
	echGroup.API.Use(jwtMiddFn)
}

func (cv *CustomValidator) CustomValidation() {
	validator := validator.New()
	_ = validator.RegisterValidation("isRequiredWith", cv.isRequiredWith)
	_ = validator.RegisterValidation("dateString", cv.dateString)
	cv.Validator = validator
}

func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.Validator.Struct(i)
}

func (cv *CustomValidator) isRequiredWith(fl validator.FieldLevel) bool {
	field := fl.Field()
	otherField, _, _ := fl.GetStructFieldOK()

	if otherField.IsValid() && otherField.Interface() != reflect.Zero(otherField.Type()).Interface() {
		if field.IsValid() && field.Interface() == reflect.Zero(field.Type()).Interface() {
			return false
		}
	}

	return true
}

func (cv *CustomValidator) dateString(fl validator.FieldLevel) bool {
	field := fl.Field()
	vValue := fl.Param()

	if vValue == "" {
		vValue = models.DateFormat
	}

	if field.Interface() == reflect.Zero(field.Type()).Interface() {
		return true
	}

	date, err := time.Parse(vValue, field.Interface().(string))

	if err != nil || (date == time.Time{}) {
		return false
	}

	return true
}

func bodyParser(c echo.Context, pl *[]byte) {
	if string(*pl) == "" {
		rawQuery := c.Request().URL.RawQuery
		m, err := url.ParseQuery(rawQuery)

		if err != nil {
			logger.Make(nil, nil).Fatal(err)
		}

		*pl, err = json.Marshal(m)

		if err != nil {
			logger.Make(nil, nil).Fatal(err)
		}
	}
}
