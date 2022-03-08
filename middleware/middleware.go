package middleware

import (
	"encoding/json"
	"gade/srv-gade-point/logger"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/referrals"

	"net/http"
	"net/url"
	"os"
	"reflect"
	"time"

	"gopkg.in/go-playground/validator.v9"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

type customMiddleware struct {
	e *echo.Echo
}

var echGroup models.EchoGroup

// InitMiddleware to generate all middleware that app need
func InitMiddleware(ech *echo.Echo, echoGroup models.EchoGroup) {
	cm := &customMiddleware{ech}
	echGroup = echoGroup
	ech.Use(middleware.RequestIDWithConfig(middleware.DefaultRequestIDConfig))
	cm.customBodyDump()

	ech.Use(middleware.Recover())
	cm.cors()
	cm.basicAuth()
	cm.jwtAuth()
	cm.customValidation()
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

func (cm *customMiddleware) customValidation() {
	validator := validator.New()
	customValidator := customValidator{}
	_ = validator.RegisterValidation("isRequiredWith", customValidator.isRequiredWith)
	_ = validator.RegisterValidation("dateString", customValidator.dateString)
	customValidator.validator = validator
	cm.e.Validator = &customValidator
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

func ReferralAuth(referralsUseCase referrals.UseCase) {
	echGroup.Referral.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			cif := c.QueryParam("cif")
			respErrors := &models.ResponseErrors{}
			response := models.Response{}
			_ = c.Bind(&cif)

			referralCIF, err := referralsUseCase.UReferralCIFValidate(c, cif)

			if err != nil {
				respErrors.SetTitle(err.Error())
				response.SetResponse("", respErrors)

				return c.JSON(http.StatusBadRequest, response)
			}

			response.Status = models.StatusSuccess
			response.Message = models.MessageDataFound
			response.Data = referralCIF

			return c.JSON(http.StatusFound, response)
		}
	})
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

func (cv *customValidator) isRequiredWith(fl validator.FieldLevel) bool {
	field := fl.Field()
	otherField, _, _ := fl.GetStructFieldOK()

	if otherField.IsValid() && otherField.Interface() != reflect.Zero(otherField.Type()).Interface() {
		if field.IsValid() && field.Interface() == reflect.Zero(field.Type()).Interface() {
			return false
		}
	}

	return true
}

func (cv *customValidator) dateString(fl validator.FieldLevel) bool {
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
