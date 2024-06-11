package cliclient


import "fmt"

// ServerError означает, что ошибка была получена при обращении к серверной части
type ServerError struct {
	Err error
}

// NewServerError конструктор
func NewServerError(err error) *ServerError {
	return &ServerError{
		Err: err,
	}
}

// Error релизует интерфейс error
func (e *ServerError) Error() string {
	return fmt.Sprintf("server error: %s", e.Err)
}

// Unwrap для errors.Unwrap
func (e *ServerError) Unwrap() error {
	return e.Err
}

// TokenNotFoundError
type TokenNotFoundError struct {
	Err error
}

// NewTokenNotFoundError конструктор
func NewTokenNotFoundError(err error) *TokenNotFoundError {
	return &TokenNotFoundError{
		Err: err,
	}
}

// Error релизует интерфейс error
func (e *TokenNotFoundError) Error() string {
	return fmt.Sprintf("token not found: %s", e.Err)
}

// Unwrap для errors.Unwrap
func (e *TokenNotFoundError) Unwrap() error {
	return e.Err
}

