package auth

import (
	"context"
	"time"

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
		refreshToken string,
		newRefreshToken string,
		newExpiredAt time.Time) (*smodel.UpdateRefreshTokenResponseModel, error)

	GetRefreshToken(
		ctx context.Context,
		request *smodel.GetRefreshTokenRequestModel) (*smodel.GetRefreshTokenResponseModel, error)

	RemoveRefreshToken(ctx context.Context, userID string, tokenID string) error
}
