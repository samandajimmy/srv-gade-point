package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/tokens"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

var response models.Response

// TokensHandler represent the httphandler for tokens
type TokensHandler struct {
	TokenUseCase tokens.UseCase
}

// NewTokensHandler represent to register tokens endpoint
func NewTokensHandler(echoGroup models.EchoGroup, tknUseCase tokens.UseCase) {
	handler := &TokensHandler{
		TokenUseCase: tknUseCase,
	}

	echoGroup.Token.POST("/create", handler.createToken)
	echoGroup.Token.GET("/get", handler.getToken)
	echoGroup.Token.GET("/refresh", handler.refreshToken)
}

func (tkn *TokensHandler) createToken(echTx echo.Context) error {
	var accountToken models.AccountToken
	response = models.Response{}
	ctx := echTx.Request().Context()
	err := echTx.Bind(&accountToken)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(getStatusCode(err), response)
	}

	err = tkn.TokenUseCase.CreateToken(ctx, &accountToken)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(getStatusCode(err), response)
	}

	accountToken.Password = ""
	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = accountToken

	return echTx.JSON(getStatusCode(err), response)
}

func (tkn *TokensHandler) getToken(echTx echo.Context) error {
	response = models.Response{}
	ctx := echTx.Request().Context()
	username := echTx.QueryParam("username")
	password := echTx.QueryParam("password")
	accToken, err := tkn.TokenUseCase.GetToken(ctx, username, password)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = accToken

	return echTx.JSON(getStatusCode(err), response)
}

func (tkn *TokensHandler) refreshToken(echTx echo.Context) error {
	response = models.Response{}
	ctx := echTx.Request().Context()
	username := echTx.QueryParam("username")
	password := echTx.QueryParam("password")
	accToken, err := tkn.TokenUseCase.RefreshToken(ctx, username, password)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(getStatusCode(err), response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = accToken

	return echTx.JSON(getStatusCode(err), response)
}

func getStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	if strings.Contains(err.Error(), "400") {
		return http.StatusBadRequest
	}

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
