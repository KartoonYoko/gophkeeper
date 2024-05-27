package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KartoonYoko/gophkeeper/internal/storage/common"
	serror "github.com/KartoonYoko/gophkeeper/internal/storage/error/auth"
	model "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) CreateUserAndRefreshToken(
	ctx context.Context,
	login string,
	password string,
	refreshTokenDurationMinute int) (*model.CreateUserAndRefreshTokenResponseModel, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Commit(ctx)

	userID := uuid.New().String()
	secretKey, err := common.GenerateSecretKey()
	if err != nil {
		return nil, fmt.Errorf("secret key generation error: %w", err)
	}

	query := `
	INSERT INTO users(id, login, password, secret_key) 
	VALUES ($1, $2, $3, $4);
	`
	_, err = tx.Exec(ctx, query, userID, login, password, secretKey)
	if err != nil {
		var pgErr *pgconn.PgError
		// already exists
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return nil, serror.NewLoginAlreadyExistsError(login)
		}

		return nil, fmt.Errorf("unable to create user: %w", err)
	}

	tokenID, err := common.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("refresh token generation error: %w", err)
	}
	query = `
	INSERT INTO user_refresh_token(token_id, user_id, expired_at)
	VALUES($1, $2, $3)
	`
	duration := time.Minute * time.Duration(refreshTokenDurationMinute)
	expiredAt := time.Now().UTC().Add(duration)
	_, err = tx.Exec(ctx, query, tokenID, userID, expiredAt)
	if err != nil {
		return nil, fmt.Errorf("unable to create user: %w", err)
	}

	result := new(model.CreateUserAndRefreshTokenResponseModel)
	result.UserID = userID
	result.ExpiredAt = expiredAt
	result.Token = tokenID

	return result, nil
}

func (s *Storage) Login(
	ctx context.Context,
	login string,
	password string,
	refreshTokenDurationMinute int) (*model.LoginResponseModel, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to start transaction: %w", err)
	}
	defer tx.Commit(ctx)

	var userID, secretKey string
	query := `SELECT id, secret_key FROM users WHERE login = $1 AND password = $2`
	err = tx.QueryRow(ctx, query, login, password).Scan(&userID, &secretKey)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, serror.NewLoginOrPasswordNotFoundError(login, password)
		}

		return nil, fmt.Errorf("unable to get user: %w", err)
	}

	tokenID, err := common.GenerateRefreshToken()
	if err != nil {
		return nil, fmt.Errorf("refresh token generation error: %w", err)
	}
	query = `
		INSERT INTO user_refresh_token(token_id, user_id, expired_at)
		VALUES($1, $2, $3)
		`
	duration := time.Minute * time.Duration(refreshTokenDurationMinute)
	expiredAt := time.Now().UTC().Add(duration)
	_, err = tx.Exec(ctx, query, tokenID, userID, expiredAt)
	if err != nil {
		return nil, fmt.Errorf("unable to create user: %w", err)
	}

	result := new(model.LoginResponseModel)
	result.UserID = userID
	result.ExpiredAt = expiredAt
	result.Token = tokenID

	return result, nil
}
