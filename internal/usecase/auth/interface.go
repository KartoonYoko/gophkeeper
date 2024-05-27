package auth

import (
	"context"

	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
)

type Storager interface {
	CreateUserAndRefreshToken(
		ctx context.Context,
		login string,
		password string,
		refreshTokenDurationMinute int) (*smodel.CreateUserAndRefreshTokenResponseModel, error)

	Login(
		ctx context.Context,
		login string,
		password string,
		refreshTokenDurationMinute int) (*smodel.LoginResponseModel, error)

	UpdateRefreshToken(
		ctx context.Context,
		userID string,
		refreshToken string,
		refreshTokenDurationMinute int) (*smodel.UpdateRefreshTokenResponseModel, error)

	RemoveRefreshToken(ctx context.Context, userID string, tokenID string) error
}
