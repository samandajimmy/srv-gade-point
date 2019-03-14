package http

import (
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/users"
	"net/http"

	"github.com/labstack/echo"
)

var response models.Response

// UserHandler represent the httphandler for user
type UserHandler struct {
	UserUseCase users.UseCase
}

// NewUserHandler represent to register user endpoint
func NewUserHandler(echoGroup models.EchoGroup, us users.UseCase) {
	handler := &UserHandler{
		UserUseCase: us,
	}

	//End Point For CMS
	echoGroup.Admin.POST("/users", handler.createUser)
	echoGroup.Admin.POST("/users/auth", handler.validateUser)
}

func (usr *UserHandler) createUser(echTx echo.Context) error {
	var user models.User
	response = models.Response{}
	ctx := echTx.Request().Context()
	err := echTx.Bind(&user)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	err = usr.UserUseCase.CreateUser(ctx, &user)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(http.StatusBadRequest, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = user

	return echTx.JSON(http.StatusOK, response)
}

// ValidateUser to get list of user and validate password
func (usr *UserHandler) validateUser(echTx echo.Context) error {
	var user models.User
	response = models.Response{}
	ctx := echTx.Request().Context()
	err := echTx.Bind(&user)

	if err != nil {
		response.Status = models.StatusError
		response.Message = models.MessageUnprocessableEntity

		return echTx.JSON(http.StatusUnprocessableEntity, response)
	}

	data, err := usr.UserUseCase.ValidateUser(ctx, user)

	if err != nil {
		response.Status = models.StatusError
		response.Message = err.Error()

		return echTx.JSON(http.StatusBadRequest, response)
	}

	response.Status = models.StatusSuccess
	response.Message = models.MessageDataSuccess
	response.Data = data

	return echTx.JSON(http.StatusOK, response)
}
