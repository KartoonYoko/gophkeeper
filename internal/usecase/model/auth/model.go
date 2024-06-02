package auth

type RegisterResponseModel struct {
	UserID       string
	RefreshToken string
}

type LoginResponseModel struct {
	UserID       string
	RefreshToken string
}

type RefreshTokenResponseModel struct {
	UserID       string
	RefreshToken string
}
