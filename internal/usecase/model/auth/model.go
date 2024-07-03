package auth

type RegisterResponseModel struct {
	UserID       string
	RefreshToken string
	SecretKey    string
}

type LoginResponseModel struct {
	UserID       string
	RefreshToken string
	SecretKey    string
}

type RefreshTokenResponseModel struct {
	UserID       string
	RefreshToken string
}
