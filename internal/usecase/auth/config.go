package auth

// Config конфигурация аутентификации
type Config struct {
	RefreshTokenDurationMinute int
	SecretJWTKey               string
	JWTDurationMinute          int
	SecretKeySecure            string
}
