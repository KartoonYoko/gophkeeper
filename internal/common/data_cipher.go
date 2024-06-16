package common

import (
	"crypto/aes"
	"crypto/cipher"
)

// DataCipherHandler реализует шифрование данных
type DataCipherHandler struct {
	aeskey []byte
	nonce  []byte

	aesblock cipher.Block
	aesgcm   cipher.AEAD
}

// NewDataCipherHandler создает шифратор данных с указанным ключом
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

// Encrypt шифрует данные
func (h *DataCipherHandler) Encrypt(data []byte) []byte {
	return h.aesgcm.Seal(nil, h.nonce, data, nil)
}

// Decrypt дешифрует данные
func (h *DataCipherHandler) Decrypt(data []byte) (encryptedname []byte, err error) {
	return h.aesgcm.Open(nil, h.nonce, data, nil)
}
