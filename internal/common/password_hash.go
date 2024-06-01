package common

import (
	"golang.org/x/crypto/bcrypt"
)

type SHA256PasswordHasher struct {
	salt string
}

func NewSHA256PasswordHasher(salt string) *SHA256PasswordHasher {
	return &SHA256PasswordHasher{salt: salt}
}

func (h *SHA256PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(h.salt+password), bcrypt.MinCost)
	return string(bytes), err
}

func (h *SHA256PasswordHasher) CheckHash(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(h.salt+password)) == nil
}
