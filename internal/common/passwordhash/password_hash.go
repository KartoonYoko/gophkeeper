/*
Package passwordhash хеширует пароль
*/
package passwordhash

import (
	"golang.org/x/crypto/bcrypt"
)

// Hasher хэширует и сравнивает пароли
type Hasher struct {
}

// New конструктор для SHA256PasswordHasher
func New() *Hasher {
	return &Hasher{}
}

// Hash возвращает хэш пароля
func (h *Hasher) Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.MinCost)
	return string(bytes), err
}

// CheckHash проверяет хэш пароля
func (h *Hasher) CheckHash(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
