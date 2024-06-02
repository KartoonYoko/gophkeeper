package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	appcommon "github.com/KartoonYoko/gophkeeper/internal/common"
	"github.com/KartoonYoko/gophkeeper/internal/logger"
	"github.com/KartoonYoko/gophkeeper/internal/storage/common"
	serror "github.com/KartoonYoko/gophkeeper/internal/storage/error/auth"
	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
	uccommon "github.com/KartoonYoko/gophkeeper/internal/usecase/common"
	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/auth"
	"go.uber.org/zap"
)

type Usecase struct {
	storage Storager
	conf    Config

	pswdHasher *appcommon.SHA256PasswordHasher
}

func New(storage Storager, config Config) *Usecase {
	uc := new(Usecase)

	uc.conf = config
	uc.storage = storage
	uc.pswdHasher = appcommon.NewSHA256PasswordHasher(uc.conf.PasswordSault)

	return uc
}

func (uc *Usecase) Register(ctx context.Context, login string, password string) (*model.RegisterResponseModel, error) {
	hpswd, err := uc.encodePassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash passwd: %w", err)
	}
	m, err := uc.storage.CreateUserAndRefreshToken(ctx, login, hpswd, uc.conf.RefreshTokenDurationMinute)
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
	// TODO
	// получить пользователя по логину (если не найдено, то serror.LoginOrPasswordNotFoundError)
	// провалидировать пароль с помощью uc.pswdHasher.CheckHash()
	// если невалидный, то вернуть ошибку serror.LoginOrPasswordNotFoundError
	// иначе создать рефреш токен для пользователя

	getUserResponse, err := uc.storage.GetUserByLogin(ctx, login)
	if err != nil {
		var exsterror *serror.LoginNotFoundError
		if errors.As(err, &exsterror) {
			return nil, serror.NewLoginNotFoundError(exsterror.Login)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// TODO работает ли?
	// проверка пароля 
	hpswd, err := uc.encodePassword(password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash passwd: %w", err)
	}
	if !uc.pswdHasher.CheckHash(getUserResponse.Password, hpswd) {
		return nil, serror.NewLoginOrPasswordNotFoundError(login, password)
	}
	
	// создать рефреш токен
	tokenID, err := uccommon.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	addRefreshTokenRequest := &smodel.CreateRefreshTokenRequestModel{
		UserID: getUserResponse.UserID,
		TokenID: tokenID,
		ExpiredAt: uc.refreshTokenExpiredAt(),
	}
	m, err := uc.storage.CreateRefreshToken(ctx, addRefreshTokenRequest)

	// resModel := new(model.LoginResponseModel)
	// resModel.UserID = 
	// resModel.RefreshToken = m.Token

	return resModel, nil
}

func (uc *Usecase) Logout(ctx context.Context, userID string, tokenID string) error {
	err := uc.storage.RemoveRefreshToken(ctx, userID, tokenID)
	if err != nil {
		return fmt.Errorf("failed to logout: %w", err)
	}

	return nil
}

func (uc *Usecase) RefreshToken(ctx context.Context, refreshToken string) (*model.RefreshTokenResponseModel, error) {
	// - получить refresh token
	// - проверить его время жизни
	// - если невалидный, то вернуть ошибку
	// - иначе обновить токен и время жизни

	r := &smodel.GetRefreshTokenRequestModel{
		TokenID: refreshToken,
	}
	tkn, err := uc.storage.GetRefreshToken(ctx, r)
	if err != nil {
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// проверить время жизни токена
	if tkn.ExpiredAt.After(time.Now().UTC()) {
		return nil, fmt.Errorf("refresh token expired")
	}

	newTkn, err := uccommon.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	tknExpiredAt := uc.refreshTokenExpiredAt()

	utkn, err := uc.storage.UpdateRefreshToken(ctx, tkn.TokenID, newTkn, tknExpiredAt)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	result := new(model.RefreshTokenResponseModel)
	result.RefreshToken = utkn.Token
	result.UserID = tkn.UserID

	return result, nil
}

func (uc *Usecase) BuildJWTString(userID string) (string, error) {
	builder := uccommon.NewJWTStringBuilder(
		uc.conf.SecretJWTKey,
		uccommon.WithUserID(userID),
		uccommon.WithTokeExpiredAtInMinute(uc.conf.JWTDurationMinute))

	return builder.BuildJWTString()
}

func (uc *Usecase) ValidateJWTString(token string) (string, error) {
	validator := uccommon.NewJWTStringValidator(uc.conf.SecretJWTKey)

	return validator.ValidateAndGetUserID(token)
}

func (uc *Usecase) encodePassword(password string) (string, error) {
	return uc.pswdHasher.Hash(password)
}

func (uc *Usecase) refreshTokenExpiredAt() time.Time {
	duration := time.Minute * time.Duration(uc.conf.RefreshTokenDurationMinute)
	return time.Now().UTC().Add(duration)
}
