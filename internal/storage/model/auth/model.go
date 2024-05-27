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
	UserID    string
	Token     string
	ExpiredAt time.Time
}