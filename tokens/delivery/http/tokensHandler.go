package http

import (
	"gade/srv-gade-point/models"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"golang.org/x/crypto/bcrypt"
)

// Response represent the response
var response = models.Response{}

// TokensHandler represent the httphandler for tokens
type TokensHandler struct {
	// nothing
}

// NewTokensHandler represent to register tokens endpoint
func NewTokensHandler(echoGroup models.EchoGroup) {
	handler := &TokensHandler{}

	echoGroup.Token.POST("/create", handler.createToken)
	echoGroup.Token.GET("/get", handler.getToken)
}

func (tkn *TokensHandler) createToken(echTx echo.Context) error {
	var accountToken models.AccountToken
	now := time.Now()
	cookieExp := stringToDuration(os.Getenv(`JWT_COOKIE_EXP`))
	tokenExp := stringToDuration(os.Getenv(`JWT_TOKEN_EXP`))
	err := echTx.Bind(&accountToken)

	if err != nil {
		return echTx.JSON(http.StatusUnprocessableEntity, accountToken)
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(accountToken.Password), bcrypt.DefaultCost)
	accountToken.Password = string(hashedPassword)

	claims := models.Token{
		accountToken.Username,
		jwt.StandardClaims{
			Id:        accountToken.Username + ":" + now.Format(models.MicroTimeFormat),
			ExpiresAt: time.Now().Add(tokenExp * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err := rawToken.SignedString([]byte(os.Getenv(`JWT_SECRET`)))

	cookie := &http.Cookie{}
	cookie.Name = accountToken.Username
	cookie.Value = token + ":" + accountToken.Password
	cookie.Expires = time.Now().Add(cookieExp * time.Hour)
	echTx.SetCookie(cookie)

	return echTx.JSON(http.StatusOK, map[string]string{
		"message":  "You were logged in!",
		"username": accountToken.Username,
		"token":    token,
	})
}

func (tkn *TokensHandler) getToken(echTx echo.Context) error {
	username := echTx.QueryParam("username")
	password := echTx.QueryParam("password")

	cookie, err := echTx.Cookie(username)

	if err != nil {
		log.Println("Your username or password were wrong!", err)
		return echTx.String(http.StatusUnauthorized, "Your username or password were wrong")
	}

	cookies := strings.Split(cookie.Value, ":")
	token, hashedPassword := cookies[0], cookies[1]
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return echTx.String(http.StatusInternalServerError, "Password does not match!")
	}

	return echTx.JSON(http.StatusOK, map[string]string{
		"message": "You're token is here!",
		"token":   token,
	})
}

func stringToDuration(str string) time.Duration {
	hours, _ := strconv.Atoi(str)

	return time.Duration(hours)
}
