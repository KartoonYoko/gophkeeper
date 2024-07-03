package store

type Config struct {
	SecretKeySecure string // секрет для шифровки/дешифровки серкетного ключа пользователя
	DataSecretKey   string // ключ для шировки/дешифровки данных
}
