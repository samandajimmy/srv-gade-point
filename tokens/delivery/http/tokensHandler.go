package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/tokens"
	"net/http"

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

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	err = tkn.TokenUseCase.CreateToken(ctx, &accountToken)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(http.StatusBadRequest, response)
	}

	accountToken.Password = ""
	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = accountToken

	return echTx.JSON(http.StatusOK, response)
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

		return echTx.JSON(http.StatusBadRequest, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = accToken

	return echTx.JSON(http.StatusOK, response)
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

		return echTx.JSON(http.StatusBadRequest, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = accToken

	return echTx.JSON(http.StatusOK, response)
}
