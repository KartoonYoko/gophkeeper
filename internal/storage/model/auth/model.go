package auth

import "time"

type CreateUserAndRefreshTokenResponseModel struct {
	UserID    string
	Token     string
	ExpiredAt time.Time
}
