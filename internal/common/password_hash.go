package common

import (
	"golang.org/x/crypto/bcrypt"
)

type SHA256PasswordHasher struct {
}

func NewSHA256PasswordHasher() *SHA256PasswordHasher {
	return &SHA256PasswordHasher{}
}

func (h *SHA256PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

func (h *SHA256PasswordHasher) CheckHash(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
