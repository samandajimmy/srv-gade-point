package usecase

import (
	"context"
	"gade/srv-gade-point/models"
	"gade/srv-gade-point/tokens"
	"time"

	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type tokenUseCase struct {
	tokenRepo      tokens.Repository
	contextTimeout time.Duration
}

// NewTokenUseCase will create new an TokenUseCase object representation of Tokens.UseCase interface
func NewTokenUseCase(tkn tokens.Repository, timeout time.Duration) tokens.UseCase {
	return &tokenUseCase{
		tokenRepo:      tkn,
		contextTimeout: timeout,
	}
}

func (tkn *tokenUseCase) CreateToken(ctx context.Context, accToken *models.AccountToken) error {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(accToken.Password), bcrypt.DefaultCost)
	accToken.Password = string(hashedPassword)
	err := tkn.tokenRepo.Create(ctx, accToken)

	if err != nil {
		log.Error(err)

		return err
	}

	return nil
}

func (tkn *tokenUseCase) GetToken(ctx context.Context, username string, password string) (*models.AccountToken, error) {
	accToken := &models.AccountToken{}
	accToken.Username = username

	// get account
	err := tkn.tokenRepo.GetByUsername(ctx, accToken)

	if err != nil {
		log.Error(err)

		return nil, models.ErrUsername
	}

	if err = verifyToken(accToken, password, false); err != nil {
		log.Error(err)

		return nil, err
	}

	// rearrange accountToken
	accToken.ID = 0
	accToken.Password = ""
	accToken.Status = nil

	return accToken, nil
}

func (tkn *tokenUseCase) RefreshToken(ctx context.Context, username string, password string) (*models.AccountToken, error) {
	accToken := &models.AccountToken{}
	accToken.Username = username

	// get account
	err := tkn.tokenRepo.GetByUsername(ctx, accToken)

	if err != nil {
		log.Error(err)

		return nil, models.ErrUsername
	}

	if err = verifyToken(accToken, password, true); err != nil {
		log.Error(err)

		return nil, err
	}

	// refresh JWT
	err = tkn.tokenRepo.UpdateToken(ctx, accToken)

	if err != nil {
		return nil, err
	}

	_ = tkn.tokenRepo.GetByUsername(ctx, accToken)

	// rearrange accountToken
	accToken.ID = 0
	accToken.Password = ""
	accToken.Status = nil

	return accToken, nil
}

func verifyToken(accToken *models.AccountToken, password string, isUpdate bool) error {
	now := time.Now()
	// validate account
	// check password
	err := bcrypt.CompareHashAndPassword([]byte(accToken.Password), []byte(password))

	if err != nil {
		return models.ErrPassword
	}

	// token availabilty
	if accToken.ExpireAt.Before(now) && !isUpdate {
		return models.ErrTokenExpired
	}

	return nil
}
