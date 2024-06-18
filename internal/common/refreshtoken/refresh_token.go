/*
Package refreshtoken реализует создание токена обновленния
*/
package refreshtoken

import (
	"crypto/rand"
	"encoding/base64"
)

// Generate генерирует рефреш токен
func Generate() (string, error) {
	n := 16
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ``, err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
