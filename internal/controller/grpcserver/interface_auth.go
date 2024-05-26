package grpcserver

import (
	"context"

	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/auth"
)

// AuthUsecase usecase для ручек авторизации/аутентификации
type AuthUsecase interface {
	// LoginUser() (secret_key string, err error)
	Register(ctx context.Context, login string, password string) (*model.RegisterResponseModel, error)
}
