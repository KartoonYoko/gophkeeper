package jwtbuilder

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/KartoonYoko/gophkeeper/internal/usecase/common/jwtclaims"
)

type Builder struct {
	secretKey string

	userID                 string
	tokenExpiredAtInMinute int
}

func New(secretKey string, opts ...func(*Builder)) Builder {
	o := Builder{}
	o.secretKey = secretKey

	for _, opt := range opts {
		opt(&o)
	}
	return o
}

// WithUserID — добавляет ID пользователя в токен
func WithUserID(userID string) func(*Builder) {
	return func(o *Builder) {
		o.userID = userID
	}
}

// WithUserID — добавляет ID пользователя в токен
func WithTokeExpiredAtInMinute(minutes int) func(*Builder) {
	return func(o *Builder) {
		o.tokenExpiredAtInMinute = minutes
	}
}

// BuildJWTString создаёт токен и возвращает его в виде строки.
func (b *Builder) BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtclaims.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(b.tokenExpiredAtInMinute))),
		},
		UserID: b.userID,
	})

	tokenString, err := token.SignedString([]byte(b.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

type JWTStringValidator struct {
	secretKey string
}

func NewJWTStringValidator(secretKey string) JWTStringValidator {
	return JWTStringValidator{secretKey: secretKey}
}

// ValidateAndGetUserID валидирует токен и получает из него UserID
func (v *JWTStringValidator) ValidateAndGetUserID(tokenString string) (string, error) {
	claims := &jwtclaims.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(v.secretKey), nil
		})
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", fmt.Errorf("token is not valid")
	}

	return claims.UserID, nil
}
