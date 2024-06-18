package auth

import (
	"context"
	"time"

	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
)

type Storager interface {
	CreateUser(
		ctx context.Context,
		request *smodel.CreateUserRequestModel) (*smodel.CreateUserReqsponseModel, error)
	UpdateRefreshToken(
		ctx context.Context,
		refreshToken string,
		newRefreshToken string,
		newExpiredAt time.Time) (*smodel.UpdateRefreshTokenResponseModel, error)

	GetUserByLogin(
		ctx context.Context,
		login string) (*smodel.GetUserByLoginResponseModel, error)
	CreateRefreshToken(
		ctx context.Context,
		request *smodel.CreateRefreshTokenRequestModel) (*smodel.CreateRefreshTokenResponseModel, error)

	GetRefreshToken(
		ctx context.Context,
		request *smodel.GetRefreshTokenRequestModel) (*smodel.GetRefreshTokenResponseModel, error)

	RemoveRefreshToken(ctx context.Context, userID string, tokenID string) error
}

type PasswordHasher interface {
	Hash(password string) (string, error)
	CheckHash(password, hash string) bool
}

type SecretKeyHandler interface {
	Encrypt(secretkey string) (encryptedname string, err error)
	Decrypt(encrypted string) (encryptedname string, err error)
	GenerateEncryptedSecretKey() (string, error)
}
