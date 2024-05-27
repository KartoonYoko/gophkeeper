package auth

import (
	"context"
	"errors"
	"fmt"

	serror "github.com/KartoonYoko/gophkeeper/internal/storage/error/auth"
	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/auth"
)

type Usecase struct {
	storage Storager
	conf Config
}

func New(storage Storager, config Config) *Usecase {
	uc := new(Usecase)

	uc.conf = config
	uc.storage = storage

	return uc
}

func (uc *Usecase) Register(ctx context.Context, login string, password string) (*model.RegisterResponseModel, error) {
	m, err := uc.storage.CreateUserAndRefreshToken(ctx, login, password, uc.conf.RefreshTokenDurationMinute)
	if err != nil {
		var exsterror *serror.LoginAlreadyExistsError
		if errors.As(err, &exsterror) {
			return nil, NewLoginAlreadyExistsError(exsterror.Login)
		}
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	resModel := new(model.RegisterResponseModel)
	resModel.UserID = m.UserID
	resModel.RefreshToken = m.Token

	return resModel, nil
}

func (uc *Usecase) Login(ctx context.Context, login string, password string) (*model.LoginResponseModel, error) {
	m, err := uc.storage.Login(ctx, login, password, uc.conf.RefreshTokenDurationMinute)
	if err != nil {
		var exsterror *serror.LoginOrPasswordNotFoundError
		if errors.As(err, &exsterror) {
			return nil, NewLoginOrPasswordNotFoundError(exsterror.Login, exsterror.Password)
		}
		return nil, fmt.Errorf("failed to login: %w", err)
	}

	resModel := new(model.LoginResponseModel)
	resModel.UserID = m.UserID
	resModel.RefreshToken = m.Token

	return resModel, nil
}
