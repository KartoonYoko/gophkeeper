package jwtvalidator

import (
	"fmt"

	"github.com/KartoonYoko/gophkeeper/internal/usecase/common/jwtclaims"
	"github.com/golang-jwt/jwt/v5"
)

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
