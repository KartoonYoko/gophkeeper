package common

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims — структура утверждений, которая включает стандартные утверждения
// и одно пользовательское — UserID
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

type JWTStringBuilder struct {
	secretKey string 

	userID string
	tokenExpiredAtInMinute int
}

func NewJWTStringBuilder(secretKey string, opts ...func(*JWTStringBuilder)) JWTStringBuilder {
	o := JWTStringBuilder{}
	o.secretKey = secretKey

    for _, opt := range opts {
        opt(&o)
    }
    return o
}

// WithUserID — добавляет ID пользователя в токен
func WithUserID(userID string) func(*JWTStringBuilder) {
	return func(o *JWTStringBuilder) {
		o.userID = userID
	}
}

// WithUserID — добавляет ID пользователя в токен
func WithTokeExpiredAtInMinute(minutes int) func(*JWTStringBuilder) {
	return func(o *JWTStringBuilder) {
		o.tokenExpiredAtInMinute = minutes
	}
}


// BuildJWTString создаёт токен и возвращает его в виде строки.
func (b *JWTStringBuilder) BuildJWTString() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
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
	claims := &Claims{}
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