package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KartoonYoko/gophkeeper/internal/common/refreshtoken"
	serror "github.com/KartoonYoko/gophkeeper/internal/storage/error/auth"
	smodel "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/common/jwtbuilder"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/common/jwtvalidator"
	model "github.com/KartoonYoko/gophkeeper/internal/usecase/model/auth"
	"github.com/google/uuid"
)

type Usecase struct {
	// Storage хранилище данных; пришлось сделать экспортируемым для моков
	Storage Storager
	conf    Config

	pswdHasher       PasswordHasher
	secretkeyHandler SecretKeyHandler
}

func New(storage Storager, pswdHasher PasswordHasher, secretkeyHandler SecretKeyHandler, config Config) *Usecase {
	uc := new(Usecase)

	uc.conf = config
	uc.Storage = storage
	uc.pswdHasher = pswdHasher
	uc.secretkeyHandler = secretkeyHandler

	return uc
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
	_, err = uc.Storage.CreateUser(ctx, createUserRequest)
	if err != nil {
		var exsterror *serror.LoginAlreadyExistsError
		if errors.As(err, &exsterror) {
			return nil, NewLoginAlreadyExistsError(exsterror.Login)
		}
		return nil, fmt.Errorf("failed to register: %w", err)
	}

	// создаём токен обновления
	rt, err := refreshtoken.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	createRefreshTokenRequest := &smodel.CreateRefreshTokenRequestModel{
		UserID:    userID,
		TokenID:   rt,
		ExpiredAt: uc.refreshTokenExpiredAt(),
	}
	_, err = uc.Storage.CreateRefreshToken(ctx, createRefreshTokenRequest)
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
		SecretKey:    secretkey,
	}

	return response, nil
}

func (uc *Usecase) Login(ctx context.Context, login string, password string) (*model.LoginResponseModel, error) {
	// получить пользователя по логину (если не найдено, то serror.LoginOrPasswordNotFoundError)
	// провалидировать пароль с помощью uc.pswdHasher.CheckHash()
	// если невалидный, то вернуть ошибку serror.LoginOrPasswordNotFoundError
	// иначе создать рефреш токен для пользователя

	getUserResponse, err := uc.Storage.GetUserByLogin(ctx, login)
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
	tokenID, err := refreshtoken.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	addRefreshTokenRequest := &smodel.CreateRefreshTokenRequestModel{
		UserID:    getUserResponse.UserID,
		TokenID:   tokenID,
		ExpiredAt: uc.refreshTokenExpiredAt(),
	}
	_, err = uc.Storage.CreateRefreshToken(ctx, addRefreshTokenRequest)
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
	err := uc.Storage.RemoveRefreshToken(ctx, userID, tokenID)
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
	tkn, err := uc.Storage.GetRefreshToken(ctx, r)
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
		return nil, NewRefreshTokenExpiredError(refreshToken, tkn.ExpiredAt)
	}

	newTkn, err := refreshtoken.Generate()
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}
	tknExpiredAt := uc.refreshTokenExpiredAt()

	utkn, err := uc.Storage.UpdateRefreshToken(ctx, tkn.TokenID, newTkn, tknExpiredAt)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token: %w", err)
	}

	result := new(model.RefreshTokenResponseModel)
	result.RefreshToken = utkn.Token
	result.UserID = tkn.UserID

	return result, nil
}

func (uc *Usecase) BuildJWTString(userID string) (string, error) {
	builder := jwtbuilder.New(
		uc.conf.SecretJWTKey,
		jwtbuilder.WithUserID(userID),
		jwtbuilder.WithTokeExpiredAtInMinute(uc.conf.JWTDurationMinute))

	return builder.BuildJWTString()
}

func (uc *Usecase) ValidateJWTString(token string) (string, error) {
	validator := jwtvalidator.NewJWTStringValidator(uc.conf.SecretJWTKey)

	return validator.ValidateAndGetUserID(token)
}

func (uc *Usecase) refreshTokenExpiredAt() time.Time {
	duration := time.Minute * time.Duration(uc.conf.RefreshTokenDurationMinute)
	return time.Now().UTC().Add(duration)
}
