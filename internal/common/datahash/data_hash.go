/*
Package datahash реализует хеширование данных
*/
package datahash

import "crypto/sha256"

type DataHasherSHA256 struct {
}

func NewDataHasherSHA256() *DataHasherSHA256 {
	return &DataHasherSHA256{}
}

func (h *DataHasherSHA256) Hash(data []byte) []byte {
	sha := sha256.New()
	sha.Write(data)
	return sha.Sum(nil)
}
