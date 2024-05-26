package auth

import (
	"context"
	"fmt"

	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/auth"
)

type Usecase struct {
	storage Storager
	conf Config
}

func New(storage Storager, config Config) *Usecase {
	uc := new(Usecase)
	uc.storage = storage

	return uc
}

func (uc *Usecase) Register(ctx context.Context, login string, password string) (*model.RegisterResponseModel, error) {
	m, err := uc.storage.CreateUserAndRefreshToken(ctx, login, password, uc.conf.RefreshTokenDurationMinute)
	if err != nil {
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	resModel := new(model.RegisterResponseModel)
	resModel.UserID = m.UserID

	return resModel, nil
}
