package auth

type Config struct {
	RefreshTokenDurationMinute int
	SecretJWTKey               string
	JWTDurationMinute          int
	PasswordSault              string
}
