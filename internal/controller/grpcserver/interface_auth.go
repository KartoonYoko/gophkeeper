package grpcserver

import (
	"context"

	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/auth"
)

// AuthUsecase usecase для ручек авторизации/аутентификации
type AuthUsecase interface {
	Register(ctx context.Context, login string, password string) (*model.RegisterResponseModel, error)
	Login(ctx context.Context, login string, password string) (*model.LoginResponseModel, error)
	Logout(ctx context.Context, userID string, tokenID string) error
	RefreshToken(ctx context.Context, refreshToken string) (*model.RefreshTokenResponseModel, error)
	BuildJWTString(userID string) (string, error)
	ValidateJWTString(token string) (string, error)
}
