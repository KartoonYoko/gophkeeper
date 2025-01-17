package auth

import "time"

type CreateUserAndRefreshTokenResponseModel struct {
	UserID    string
	Token     string
	ExpiredAt time.Time
}

type LoginResponseModel struct {
	UserID    string
	Token     string
	ExpiredAt time.Time
}

type UpdateRefreshTokenResponseModel struct {
	Token     string
	ExpiredAt time.Time
}

type GetRefreshTokenRequestModel struct {
	TokenID string
}

type GetRefreshTokenResponseModel struct {
	TokenID   string
	UserID    string
	ExpiredAt time.Time
}

type GetUserByLoginResponseModel struct {
	UserID    string
	Password  string
	SecretKey string
}

type CreateRefreshTokenRequestModel struct {
	UserID    string
	TokenID   string
	ExpiredAt time.Time
}

type CreateRefreshTokenResponseModel struct {
}

type CreateUserRequestModel struct {
	UserID    string
	Login     string
	SecretKey string
	Password  string
}

type CreateUserReqsponseModel struct {
}
