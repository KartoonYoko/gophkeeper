package common

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
)

type SecretKeyHandler struct {
	aeskey []byte
	nonce  []byte

	aesblock cipher.Block
	aesgcm   cipher.AEAD
}

func NewSecretKeyHandler(key string) (*SecretKeyHandler, error) {
	keysize := 2 * aes.BlockSize

	aeskey := getSliceNFromString(key, keysize)

	aesblock, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	// todo: возможно ли использовать при шифровке/дешифровке каждый раз разный nonce?
	nonce := getSliceNFromString(key, aesgcm.NonceSize())

	h := new(SecretKeyHandler)
	h.aeskey = aeskey
	h.nonce = nonce
	h.aesblock = aesblock
	h.aesgcm = aesgcm

	return h, nil
}

func (h *SecretKeyHandler) Encrypt(secretkey string) (encryptedname string, err error) {
	dst := h.aesgcm.Seal(nil, h.nonce, []byte(secretkey), nil)
	return hex.EncodeToString(dst), nil
}

func (h *SecretKeyHandler) Decrypt(encrypted string) (encryptedname string, err error) {
	encd, err := hex.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	dst, err := h.aesgcm.Open(nil, h.nonce, encd, nil)
	if err != nil {
		return "", err
	}

	return string(dst), nil
}

func (h *SecretKeyHandler) GenerateEncryptedSecretKey() (string, error) {
	sc, err := generateSecretKey()
	if err != nil {
		return "", err
	}

	return h.Encrypt(sc)
}

func generateSecretKey() (string, error) {
	n := 16
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return ``, err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

func getSliceNFromString(str string, n int) []byte {
	bytes := []byte(str)
	keyb := make([]byte, n)
	copy(bytes, keyb)

	return keyb
}
