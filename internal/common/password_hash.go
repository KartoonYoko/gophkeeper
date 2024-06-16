package common

import (
	"golang.org/x/crypto/bcrypt"
)

// SHA256PasswordHasher хэширует и сравнивает пароли
type SHA256PasswordHasher struct {
}

// NewSHA256PasswordHasher конструктор для SHA256PasswordHasher
func NewSHA256PasswordHasher() *SHA256PasswordHasher {
	return &SHA256PasswordHasher{}
}

// Hash возвращает хэш пароля
func (h *SHA256PasswordHasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

// CheckHash проверяет хэш пароля
func (h *SHA256PasswordHasher) CheckHash(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
