package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	serror "github.com/KartoonYoko/gophkeeper/internal/storage/error/auth"
	model "github.com/KartoonYoko/gophkeeper/internal/storage/model/auth"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Storage) GetUserByLogin(
	ctx context.Context,
	login string) (*model.GetUserByLoginResponseModel, error) {
	var userID, password, secretkey string
	query := `SELECT id, password, secret_key FROM users WHERE login = $1`
	err := s.pool.QueryRow(ctx, query, login).Scan(&userID, &password, &secretkey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, serror.NewLoginNotFoundError(login)
		}
		return nil, fmt.Errorf("unable to get user: %w", err)
	}

	res := &model.GetUserByLoginResponseModel{
		UserID:   userID,
		Password: password,
		SecretKey: secretkey,
	}

	return res, nil
}

func (s *Storage) CreateUser(
	ctx context.Context,
	request *model.CreateUserRequestModel) (*model.CreateUserReqsponseModel, error) {
	query := `
		INSERT INTO users(id, login, password, secret_key) 
		VALUES ($1, $2, $3, $4);
		`
	_, err := s.pool.Exec(ctx, query, request.UserID, request.Login, request.Password, request.SecretKey)
	if err != nil {
		var pgErr *pgconn.PgError
		// already exists
		if errors.As(err, &pgErr) && pgerrcode.UniqueViolation == pgErr.Code {
			return nil, serror.NewLoginAlreadyExistsError(request.Login)
		}

		return nil, fmt.Errorf("unable to create user: %w", err)
	}

	return &model.CreateUserReqsponseModel{}, nil
}

func (s *Storage) CreateRefreshToken(
	ctx context.Context,
	request *model.CreateRefreshTokenRequestModel) (*model.CreateRefreshTokenResponseModel, error) {
	query := `
		INSERT INTO user_refresh_token(token_id, user_id, expired_at)
		VALUES($1, $2, $3)
		`
	_, err := s.pool.Exec(ctx, query, request.TokenID, request.UserID, request.ExpiredAt)
	if err != nil {
		return nil, fmt.Errorf("unable to create refresh token: %w", err)
	}

	response := &model.CreateRefreshTokenResponseModel{}

	return response, nil
}

func (s *Storage) RemoveRefreshToken(ctx context.Context, userID string, tokenID string) error {
	query := `DELETE FROM user_refresh_token WHERE user_id = $1 AND token_id = $2`
	_, err := s.pool.Exec(ctx, query, userID, tokenID)
	if err != nil {
		return fmt.Errorf("unable to remove refresh token: %w", err)
	}

	return nil
}

func (s *Storage) GetRefreshToken(ctx context.Context, request *model.GetRefreshTokenRequestModel) (*model.GetRefreshTokenResponseModel, error) {
	var userID string
	var expiredAt time.Time
	query := `SELECT user_id, expired_at FROM user_refresh_token WHERE token_id=$1`
	err := s.pool.QueryRow(ctx, query, request.TokenID).Scan(&userID, &expiredAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, serror.NewNotFoundError(err)
		}
		return nil, fmt.Errorf("unable find refresh token: %w", err)
	}

	return &model.GetRefreshTokenResponseModel{
		TokenID:   request.TokenID,
		UserID:    userID,
		ExpiredAt: expiredAt,
	}, nil
}

func (s *Storage) UpdateRefreshToken(
	ctx context.Context,
	refreshToken string,
	newRefreshToken string,
	newExpiredAt time.Time) (*model.UpdateRefreshTokenResponseModel, error) {

	query := `UPDATE user_refresh_token SET token_id=$1, expired_at=$2 WHERE token_id=$3`
	_, err := s.pool.Exec(ctx, query, newRefreshToken, newExpiredAt, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("unable find refresh token: %w", err)
	}

	return &model.UpdateRefreshTokenResponseModel{
		Token:     newRefreshToken,
		ExpiredAt: newExpiredAt,
	}, nil
}

