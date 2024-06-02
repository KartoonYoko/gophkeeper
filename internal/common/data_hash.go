package common

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/hex"
)

type DataCipherHandler struct {
	aeskey []byte
	nonce  []byte

	aesblock cipher.Block
	aesgcm   cipher.AEAD
}

func NewDataCipherHandler(key string) (*DataCipherHandler, error) {
	keysize := 2 * aes.BlockSize

	aeskey := make([]byte, keysize)
	copy([]byte(key), aeskey)

	aesblock, err := aes.NewCipher(aeskey)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(aesblock)
	if err != nil {
		return nil, err
	}

	// todo: возможно ли использовать при шифровке/дешифровке каждый раз разный nonce?
	nonce := make([]byte, aesgcm.NonceSize())
	copy([]byte(key), nonce)

	h := new(DataCipherHandler)
	h.aeskey = aeskey
	h.nonce = nonce
	h.aesblock = aesblock
	h.aesgcm = aesgcm

	return h, nil
}

func (h *DataCipherHandler) Encrypt(data []byte) []byte {
	return h.aesgcm.Seal(nil, h.nonce, data, nil)
}

func (h *DataCipherHandler) Decrypt(encrypted string) (encryptedname string, err error) {
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
