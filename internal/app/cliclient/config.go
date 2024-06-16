package cliclient

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

// Config конфигурация клиентского приложения
type Config struct {
	ServerAddress  string // адрес сервера; флаг a
	ConfigFileName string // имя конфигурационного файла; флаг c

	wasSetServerAddress  bool
	wasSetConfigFileName bool
}

type configFileJSON struct {
	ServerAddress *string `json:"server_address"` // аналог переменной окружения SERVER_ADDRESS или флага -a
}

// NewConfig собирает конфигурацию из флагов командной строки, переменных среды
func NewConfig() (*Config, error) {
	c := new(Config)

	err := c.setFromFlags()
	if err != nil {
		return nil, err
	}

	err = c.setFromEnv()
	if err != nil {
		return nil, err
	}

	err = c.setFromConfigFile()
	if err != nil {
		return nil, err
	}

	return c, nil
}

// setFromEnv устанавливает данные из переменных окружения, если они не были заданы ранее
func (c *Config) setFromEnv() error {
	if !c.wasSetServerAddress {
		envValue, ok := os.LookupEnv("SERVER_ADDRESS")
		c.wasSetServerAddress = ok
		if ok {
			c.ServerAddress = envValue
		}
	}

	if !c.wasSetConfigFileName {
		envValue, ok := os.LookupEnv("CONFIG_FILE_NAME")
		c.wasSetConfigFileName = ok
		if ok {
			c.ConfigFileName = envValue
		}
	}

	return nil
}

// setFromFlags устанавливает данные из переданных флагов или данные по умолчанию
func (c *Config) setFromFlags() error {
	a := flag.String("a", ":8080", "Flag responsible for http server start")
	f := flag.String("c", "config.json", "Configuration file name")
	flag.Parse()

	c.ServerAddress = *a
	c.ConfigFileName = *f

	c.wasSetServerAddress = isFlagPassed("a")
	c.wasSetServerAddress = isFlagPassed("c")

	return nil
}

// setFromConfigFile устанавливает только те данные из файла, которые ещё не заданы
func (c *Config) setFromConfigFile() error {
	if c.ConfigFileName == "" {
		return nil
	}
	f, err := os.Open(c.ConfigFileName)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}
		return fmt.Errorf("can not open config file: %w", err)
	}
	defer f.Close()

	var b []byte
	b, err = io.ReadAll(f)
	if err != nil {
		return fmt.Errorf("can not read config file: %w", err)
	}

	var j configFileJSON
	err = json.Unmarshal(b, &j)
	if err != nil {
		return fmt.Errorf("can not unmarshal config file: %w", err)
	}

	if !c.wasSetServerAddress && j.ServerAddress != nil {
		c.ServerAddress = *j.ServerAddress
		c.wasSetServerAddress = true
	}
	return nil
}

// isFlagPassed определяет был ли передан флаг
func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
