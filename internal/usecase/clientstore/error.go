package clientstore

import "fmt"

// ServerError означает, что ошибка была получена при обращении к серверной части
type ServerError struct {
	Err error
}

// NewServerError конструктор
func NewServerError(err error) *ServerError {
	return &ServerError{}
}

// Error релизует интерфейс error
func (e *ServerError) Error() string {
	return fmt.Sprintf("server error: %s", e.Err)
}

// Unwrap для errors.Unwrap
func (e *ServerError) Unwrap() error {
	return e.Err
}
