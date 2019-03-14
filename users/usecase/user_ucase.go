package usecase

import (
	"context"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/users"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo       users.Repository
	contextTimeout time.Duration
}

func NewUserUseCase(usr users.Repository, timeout time.Duration) users.UseCase {
	return &userUseCase{
		userRepo:       usr,
		contextTimeout: timeout,
	}
}

func (usr *userUseCase) CreateUser(ctx context.Context, user *models.User) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)
	err := usr.userRepo.Create(ctx, user)

	if err != nil {
		log.Error(err)

		return err
	}

	return nil
}

func (usr *userUseCase) ValidateUser(ctx context.Context, userPayload models.User) (*models.User, error) {
	user := userPayload

	// get user
	err := usr.userRepo.GetByUsername(ctx, &user)
	if err != nil {
		log.Error(err)

		return nil, models.ErrUsername
	}

	if err = verifyToken(&user, userPayload.Password); err != nil {
		log.Error(err)

		return nil, err
	}

	// validate
	user.ID = 0
	user.Password = ""
	user.Status = nil

	return &user, nil
}

func verifyToken(user *models.User, password string) error {
	// validate password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return models.ErrPassword
	}

	return nil
}
