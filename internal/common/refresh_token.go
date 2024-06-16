package common

import (
	"crypto/rand"
	"encoding/base64"
)

// GenerateRefreshToken генерирует рефреш токен
func GenerateRefreshToken() (string, error) {
	n := 16
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ``, err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
