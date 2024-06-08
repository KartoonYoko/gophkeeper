package clientstorage

import (
	"context"
	"os"
	"encoding/json"
)

type Storage struct {
	rootPath string

	tokensFileName string
}

func New() (*Storage, error) {
	s := new(Storage)
	hd, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	rootPath := hd + string(os.PathSeparator) + ".gophkeeper"
	err = os.MkdirAll(rootPath, os.ModePerm)
	if err != nil {
		return nil, err
	}
	
	s.rootPath = rootPath
	s.tokensFileName = "tokens.json"

	return s, nil
}

func (s *Storage) SaveTokens(ctx context.Context, at string, rt string) error {
	b, err := json.Marshal(tokensFile{at, rt})
	if err != nil {
		return err
	}
	err = os.WriteFile(s.getTokensPath(), b, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) GetTokens() (at string, rt string, err error) {
	b, err := os.ReadFile(s.getTokensPath())
	if err != nil {
		return ``, ``, err
	}

	tf := &tokensFile{at, rt}
	err = json.Unmarshal(b, tf)
	if err != nil {
		return ``, ``, err
	}

	return tf.AccessToken, tf.RefreshToken, nil
}

func (s *Storage) RemoveTokens() error {
	return os.Remove(s.getTokensPath())
}

func (s *Storage) SaveData(ctx context.Context, fileName string, data []byte) error {
	err := os.WriteFile(s.getDataPathWithName(fileName), data, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func (s *Storage) getTokensPath() string {
	return s.rootPath + string(os.PathSeparator) + s.tokensFileName
}

func (s *Storage) getDataPathWithName(filename string) string {
	return s.rootPath + string(os.PathSeparator) + filename
}
