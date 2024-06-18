package auth

import (
	"fmt"
	"time"
)

// LoginAlreadyExistsError говорит о том, что логин уже существует в БД
type LoginAlreadyExistsError struct {
	Login string // логин пользователя, который уже существует в базе
}

// NewLoginAlreadyExistsError Конструктор для URLAlreadyExistsError
func NewLoginAlreadyExistsError(login string) *LoginAlreadyExistsError {
	return &LoginAlreadyExistsError{
		Login: login,
	}
}

// Error релизует интерфейс error
func (e *LoginAlreadyExistsError) Error() string {
	return fmt.Sprintf("login %s already exists", e.Login)
}

// LoginOrPasswordNotFoundError говорит о том, что логин или пароль не найдены
type LoginOrPasswordNotFoundError struct {
	Login    string
	Password string
}

// NewLoginOrPasswordNotFoundError конструктор
func NewLoginOrPasswordNotFoundError(login string, password string) *LoginOrPasswordNotFoundError {
	return &LoginOrPasswordNotFoundError{
		Login:    login,
		Password: password,
	}
}

// Error релизует интерфейс error
func (e *LoginOrPasswordNotFoundError) Error() string {
	return fmt.Sprintf("login %s or password not found", e.Login)
}

// LoginNotFoundError
type LoginNotFoundError struct {
	Login string
}

// NewLoginNotFoundError конструктор
func NewLoginNotFoundError(login string) *LoginNotFoundError {
	return &LoginNotFoundError{
		Login: login,
	}
}

// Error релизует интерфейс error
func (e *LoginNotFoundError) Error() string {
	return fmt.Sprintf("login %s not found", e.Login)
}

// RefreshTokenNotFoundError
type RefreshTokenNotFoundError struct {
	Token string
}

// NewRefreshTokenNotFoundError конструктор
func NewRefreshTokenNotFoundError(token string) *RefreshTokenNotFoundError {
	return &RefreshTokenNotFoundError{
		Token: token,
	}
}

// Error релизует интерфейс error
func (e *RefreshTokenNotFoundError) Error() string {
	return fmt.Sprintf("token %s not found", e.Token)
}

// RefreshTokenExpiredError
type RefreshTokenExpiredError struct {
	Token     string
	ExpiredAt time.Time
}

// NewRefreshTokenNotFoundError конструктор
func NewRefreshTokenExpiredError(token string, expiredAt time.Time) *RefreshTokenExpiredError {
	return &RefreshTokenExpiredError{
		Token:     token,
		ExpiredAt: expiredAt,
	}
}

// Error релизует интерфейс error
func (e *RefreshTokenExpiredError) Error() string {
	return fmt.Sprintf("token %s expired at %s", e.Token, e.ExpiredAt)
}
