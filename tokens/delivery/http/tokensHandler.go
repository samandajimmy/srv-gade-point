package http

import (
	"gade/srv-gade-point/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

var response models.Response

// TokensHandler represent the httphandler for tokens
type TokensHandler struct {
	// nothing
}

// NewTokensHandler represent to register tokens endpoint
func NewTokensHandler(echoGroup models.EchoGroup) {
	handler := &TokensHandler{}

	echoGroup.Token.POST("/create", handler.createToken)
	echoGroup.Token.GET("/get", handler.getToken)
	echoGroup.Token.GET("/refresh", handler.refreshToken)
}

func (tkn *TokensHandler) createToken(echTx echo.Context) error {
	var accountToken models.AccountToken
	response = models.Response{}
	now := time.Now()
	tokenExp := now.Add(stringToDuration(os.Getenv(`JWT_TOKEN_EXP`)) * time.Hour)
	err := echTx.Bind(&accountToken)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(accountToken.Password), bcrypt.DefaultCost)
	accountToken.Password = string(hashedPassword)
	token, err := createJWTToken(accountToken, now, tokenExp)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageTokenFailed

		return echTx.JSON(http.StatusServiceUnavailable, response)
	}

	cookie := &http.Cookie{}
	cookie.Name = accountToken.Username
	cookie.Value = token + ":" + accountToken.Password + ":" + tokenExp.Format(models.MicroTimeFormat)
	cookie.Expires = tokenExp
	echTx.SetCookie(cookie)

	accountToken.Password = ""
	accountToken.ExpiresAt = tokenExp.Format(models.MicroTimeFormat)

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = accountToken

	return echTx.JSON(http.StatusOK, response)
}

func (tkn *TokensHandler) getToken(echTx echo.Context) error {
	response = models.Response{}
	username := echTx.QueryParam("username")
	password := echTx.QueryParam("password")
	statusCode, msg, accountToken := verifyAccount(echTx, username, password)

	if msg != "" {
		response.Status = models.StatusError
		response.Message = msg

		return echTx.JSON(statusCode, response)
	}

	accountToken.Password = ""
	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = accountToken

	return echTx.JSON(http.StatusOK, response)
}

func (tkn *TokensHandler) refreshToken(echTx echo.Context) error {
	response = models.Response{}
	username := echTx.QueryParam("username")
	password := echTx.QueryParam("password")
	now := time.Now()
	tokenExp := now.Add(stringToDuration(os.Getenv(`JWT_TOKEN_EXP`)) * time.Hour)
	statusCode, msg, accountToken := verifyAccount(echTx, username, password)

	if msg != "" {
		response.Status = models.StatusError
		response.Message = msg

		return echTx.JSON(statusCode, response)
	}

	token, err := createJWTToken(*accountToken, now, tokenExp)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageTokenFailed

		return echTx.JSON(http.StatusServiceUnavailable, response)
	}

	cookie, _ := echTx.Cookie(username)
	cookie.Value = token + ":" + accountToken.Password + ":" + tokenExp.Format(models.MicroTimeFormat)
	cookie.Expires = tokenExp
	echTx.SetCookie(cookie)

	accountToken = &models.AccountToken{
		Username:  username,
		Token:     token,
		ExpiresAt: tokenExp.Format(models.MicroTimeFormat),
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = accountToken

	return echTx.JSON(http.StatusOK, response)
}

func verifyAccount(echTx echo.Context, username string, password string) (int, string, *models.AccountToken) {
	cookie, err := echTx.Cookie(username)

	if err != nil {
		log.Println("Your username was wrong!", err)
		return http.StatusUnauthorized, "Your username was wrong", nil
	}

	cookies := strings.Split(cookie.Value, ":")
	token, hashedPassword, tokenExp := cookies[0], cookies[1], cookies[2]
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return http.StatusUnauthorized, "Password does not match!", nil
	}

	accountToken := models.AccountToken{
		Username:  username,
		Token:     token,
		Password:  hashedPassword,
		ExpiresAt: tokenExp,
	}

	return http.StatusOK, "", &accountToken
}

func createJWTToken(accountToken models.AccountToken, now time.Time, tokenExp time.Time) (string, error) {
	claims := models.Token{
		accountToken.Username,
		jwt.StandardClaims{
			Id:        accountToken.Username,
			ExpiresAt: tokenExp.Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return rawToken.SignedString([]byte(os.Getenv(`JWT_SECRET`)))
}

func stringToDuration(str string) time.Duration {
	hours, _ := strconv.Atoi(str)

	return time.Duration(hours)
}
