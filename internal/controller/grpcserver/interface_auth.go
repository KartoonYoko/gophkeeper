package grpcserver

// AuthUsecase usecase для ручек авторизации/аутентификации
type AuthUsecase interface {
	LoginUser() (secret_key string, err error)
	RegisterUser() (secret_key string, err error)
}