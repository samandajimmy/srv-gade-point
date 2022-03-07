package usecase

import (
	"context"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/users"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo users.Repository
}

func NewUserUseCase(usr users.Repository) users.UseCase {
	return &userUseCase{
		userRepo: usr,
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

	// rearrange data
	user.Password = ""

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

	if err = verifyPassword(&user, userPayload.Password); err != nil {
		log.Error(err)

		return nil, err
	}

	// rearrange data
	user.ID = 0
	user.Password = ""
	user.Status = nil

	return &user, nil
}

func verifyPassword(user *models.User, password string) error {
	// validate password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return models.ErrPassword
	}

	return nil
}
