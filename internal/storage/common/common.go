package common

import (
	"crypto/rand"
	"encoding/base64"
)

func GenerateSecretKey() (string, error) {
	n := 16
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ``, err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func GenerateRefreshToken() (string, error) {
	n := 16
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ``, err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}