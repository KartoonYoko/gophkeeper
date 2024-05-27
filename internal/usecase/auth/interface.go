package auth

import (
	"context"

	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
)

type Storager interface {
	// Login(ctx context.Context, request *pb.LoginRequest) (*pb.LoginResponse, error)
	// Logout(ctx context.Context, request *pb.LogoutRequest) (*pb.LogoutResponse, error)
	// RefreshToken(ctx context.Context, request *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error)
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
}
