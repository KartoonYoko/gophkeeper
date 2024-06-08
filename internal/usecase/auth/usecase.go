package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	appcommon "github.com/KartoonYoko/gophkeeper/internal/common"
	serror "github.com/KartoonYoko/gophkeeper/internal/storage/error/auth"
	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
	uccommon "github.com/KartoonYoko/gophkeeper/internal/usecase/common"
	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/auth"
	"github.com/google/uuid"
)

type Usecase struct {
	storage Storager
	conf    Config

	pswdHasher       *appcommon.SHA256PasswordHasher
	secretkeyHandler *appcommon.SecretKeyHandler
}

func New(storage Storager, config Config) (*Usecase, error) {
	var err error
	uc := new(Usecase)

	uc.conf = config
	uc.storage = storage
	uc.pswdHasher = appcommon.NewSHA256PasswordHasher()
	uc.secretkeyHandler, err = appcommon.NewSecretKeyHandler(uc.conf.SecretKeySecure)
	if err != nil {
		return nil, fmt.Errorf("unable to create secret key handler: %v", err)
	}

	return uc, nil
}

func (uc *Usecase) Register(ctx context.Context, login string, password string) (*model.RegisterResponseModel, error) {
	// создаём пользователя
	sc, err := uc.secretkeyHandler.GenerateEncryptedSecretKey()
	if err != nil {
		return nil, fmt.Errorf("failed to generate secret key: %w", err)
	}
	hpswd, err := uc.pswdHasher.Hash(password)
	if err != nil {
		return nil, fmt.Errorf("failed to generate password: %w", err)
	}
	userID := uuid.New().String()
	createUserRequest := &smodel.CreateUserRequestModel{
		Login:     login,
		Password:  hpswd,
		UserID:    userID,
		SecretKey: sc,
	}
	_, err = uc.storage.CreateUser(ctx, createUserRequest)
	if err != nil {
		var exsterror *serror.LoginAlreadyExistsError
		if errors.As(err, &exsterror) {
			return nil, NewLoginAlreadyExistsError(exsterror.Login)
		}
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	// создаём токен обновления
	rt, err := appcommon.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	createRefreshTokenRequest := &smodel.CreateRefreshTokenRequestModel{
		UserID:    userID,
		TokenID:   rt,
		ExpiredAt: uc.refreshTokenExpiredAt(),
	}
	_, err = uc.storage.CreateRefreshToken(ctx, createRefreshTokenRequest)
	if err != nil {
		// todo удалить пользователя
		return nil, fmt.Errorf("failed create refresh token: %w", err)
	}

	secretkey, err := uc.secretkeyHandler.Decrypt(sc)
	if err != nil {
		return nil, fmt.Errorf("failed decrypt secret key: %w", err)
	}
	response := &model.RegisterResponseModel{
		RefreshToken: rt,
		UserID:       userID,
		SecretKey: secretkey,
	}

	return response, nil
}

func (uc *Usecase) Login(ctx context.Context, login string, password string) (*model.LoginResponseModel, error) {
	// получить пользователя по логину (если не найдено, то serror.LoginOrPasswordNotFoundError)
	// провалидировать пароль с помощью uc.pswdHasher.CheckHash()
	// если невалидный, то вернуть ошибку serror.LoginOrPasswordNotFoundError
	// иначе создать рефреш токен для пользователя

	getUserResponse, err := uc.storage.GetUserByLogin(ctx, login)
	if err != nil {
		var exsterror *serror.LoginNotFoundError
		if errors.As(err, &exsterror) {
			return nil, NewLoginNotFoundError(exsterror.Login)
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// проверка пароля
	if !uc.pswdHasher.CheckHash(password, getUserResponse.Password) {
		return nil, NewLoginOrPasswordNotFoundError(login, password)
	}

	// создать рефреш токен
	tokenID, err := uccommon.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	addRefreshTokenRequest := &smodel.CreateRefreshTokenRequestModel{
		UserID:    getUserResponse.UserID,
		TokenID:   tokenID,
		ExpiredAt: uc.refreshTokenExpiredAt(),
	}
	_, err = uc.storage.CreateRefreshToken(ctx, addRefreshTokenRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create refresh token: %w", err)
	}

	resModel := new(model.LoginResponseModel)
	resModel.UserID = getUserResponse.UserID
	resModel.RefreshToken = tokenID
	resModel.SecretKey = getUserResponse.SecretKey

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
		// если токена не существует, то вернуть ошибку
		var notfound *serror.NotFoundError
		if errors.As(err, &notfound) {
			return nil, NewRefreshTokenNotFoundError(refreshToken)
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	// проверить время жизни токена
	if tkn.ExpiredAt.Before(time.Now().UTC()) {
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

func (uc *Usecase) refreshTokenExpiredAt() time.Time {
	duration := time.Minute * time.Duration(uc.conf.RefreshTokenDurationMinute)
	return time.Now().UTC().Add(duration)
}
