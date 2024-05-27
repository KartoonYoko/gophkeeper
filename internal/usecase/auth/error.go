package auth

import "fmt"

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
	Login string
	Password string
}

// NewLoginOrPasswordNotFoundError конструктор
func NewLoginOrPasswordNotFoundError(login string, password string) *LoginOrPasswordNotFoundError {
	return &LoginOrPasswordNotFoundError{
		Login: login,
		Password: password,
	}
}

// Error релизует интерфейс error
func (e *LoginOrPasswordNotFoundError) Error() string {
	return fmt.Sprintf("login %s or password not found", e.Login)
}
