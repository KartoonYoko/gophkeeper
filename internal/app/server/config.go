package server

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

// Config - конфигурация для серверной части
type Config struct {
	ServerAddress              string // адрес сервера; флаг a
	DatabaseDsn                string // строка подключения к БД; флаг d
	MinioAddress               string // адрес minio сервера; флаг ma
	MinioAccessKeyID           string // логин minio; флаг mak
	MinioSecretAccessKey       string // пароль minio; флаг msa
	SecretJWTKey               string // ключ для подписи JWT токенов; флаг jk
	JWTTokenLifetimeMinutes    int    // время жизни JWT токенов в минутах; флаг jlt
	RefreshTokenLifeimeMinutes int    // время жизни токена обновления в минутах; флаг rlt
	UserSecretKeySecure        string // ключ для шифроваки ключей пользователей; флаг us
	DataSecretKeySecure        string // ключ для шифрования данных; флаг ds

	wasSetServerAddress              bool
	wasSetDatabaseDsn                bool
	wasSetMinioAddress               bool
	wasSetMinioAccessKeyID           bool
	wasSetMinioSecretAccessKey       bool
	wasSetSecretJWTKey               bool
	wasSetJWTTokenLifetimeMinutes    bool
	wasSetRefreshTokenLifeimeMinutes bool
	wasSetUserSecretKeySecure        bool
	wasSetDataSecretKeySecure        bool
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

	return c, nil
}

// setFromEnv устанавливает данные из переменных окружения, если они не были заданы ранее
func (c *Config) setFromEnv() error {
	if !c.wasSetServerAddress {
		envValue, ok := os.LookupEnv("SERVER_ADDRESS")
		c.wasSetDatabaseDsn = ok
		if ok {
			c.ServerAddress = envValue
		}
	}

	if !c.wasSetDatabaseDsn {
		envValue, ok := os.LookupEnv("DATABASE_DSN")
		c.wasSetDatabaseDsn = ok
		if ok {
			c.DatabaseDsn = envValue
		}
	}

	if !c.wasSetMinioAddress {
		envValue, ok := os.LookupEnv("MINIO_ADDRESS")
		c.wasSetMinioAddress = ok
		if ok {
			c.MinioAddress = envValue
		}
	}

	if !c.wasSetMinioAccessKeyID {
		envValue, ok := os.LookupEnv("MINIO_ACCESS_KEY_ID")
		c.wasSetMinioAccessKeyID = ok
		if ok {
			c.MinioAccessKeyID = envValue
		}
	}

	if !c.wasSetMinioSecretAccessKey {
		envValue, ok := os.LookupEnv("MINIO_SECRET_KEY")
		c.wasSetMinioSecretAccessKey = ok
		if ok {
			c.MinioSecretAccessKey = envValue
		}
	}

	// -----
	if !c.wasSetSecretJWTKey {
		envValue, ok := os.LookupEnv("JWT_SECRET_KEY")
		c.wasSetSecretJWTKey = ok
		if ok {
			c.SecretJWTKey = envValue
		}
	}

	if !c.wasSetJWTTokenLifetimeMinutes {
		envValue, ok := os.LookupEnv("JWT_LIFETIME_MINUTES")
		c.wasSetJWTTokenLifetimeMinutes = ok
		if ok {
			evnInt, err := strconv.Atoi(envValue)
			if err != nil {
				return fmt.Errorf("can not parse env value for jwt lifetime minutes: %w", err)
			}
			c.JWTTokenLifetimeMinutes = evnInt
		}
	}

	if !c.wasSetRefreshTokenLifeimeMinutes {
		envValue, ok := os.LookupEnv("JWT_REFRESH_LIFETIME_MINUTES")
		c.wasSetRefreshTokenLifeimeMinutes = ok
		if ok {
			evnInt, err := strconv.Atoi(envValue)
			if err != nil {
				return fmt.Errorf("can not parse env value for refresh jwt lifetime minutes: %w", err)
			}
			c.RefreshTokenLifeimeMinutes = evnInt
		}
	}

	if !c.wasSetDataSecretKeySecure {
		envValue, ok := os.LookupEnv("DATA_SECRET_KEY_SECURE")
		c.wasSetDataSecretKeySecure = ok
		if ok {
			c.DataSecretKeySecure = envValue
		}
	}

	if !c.wasSetUserSecretKeySecure {
		envValue, ok := os.LookupEnv("USER_SECRET_KEY_SECURE")
		c.wasSetUserSecretKeySecure = ok
		if ok {
			c.UserSecretKeySecure = envValue
		}
	}

	return nil
}

// setFromFlags устанавливает данные из переданных флагов или данные по умолчанию
func (c *Config) setFromFlags() error {
	a := flag.String("a", ":8080", "Flag responsible for grpc server start")
	d := flag.String("d", "", "Database connection string")
	ma := flag.String("ma", "", "Minio address")
	mak := flag.String("mak", "", "Minio login")
	msa := flag.String("msa", "", "Minio password")
	jk := flag.String("jk", "supersecret", "JWT secret key")
	jlt := flag.Int("jlt", 60, "JWT lifetime minutes")
	rlt := flag.Int("rlt", 3600, "Refresh token lifetime minutes")
	us := flag.String("us", "default", "Key to encode user secret key")
	ds := flag.String("ds", "default", "Key to encode data")

	flag.Parse()

	c.ServerAddress = *a
	c.DatabaseDsn = *d
	c.MinioAddress = *ma
	c.MinioAccessKeyID = *mak
	c.MinioSecretAccessKey = *msa
	c.SecretJWTKey = *jk
	c.JWTTokenLifetimeMinutes = *jlt
	c.RefreshTokenLifeimeMinutes = *rlt
	c.UserSecretKeySecure = *us
	c.DataSecretKeySecure = *ds

	c.wasSetServerAddress = isFlagPassed("a")
	c.wasSetDatabaseDsn = isFlagPassed("d")
	c.wasSetMinioAddress = isFlagPassed("ma")
	c.wasSetMinioAccessKeyID = isFlagPassed("mak")
	c.wasSetMinioSecretAccessKey = isFlagPassed("msa")
	c.wasSetSecretJWTKey = isFlagPassed("jk")
	c.wasSetJWTTokenLifetimeMinutes = isFlagPassed("jlt")
	c.wasSetRefreshTokenLifeimeMinutes = isFlagPassed("rlt")
	c.wasSetUserSecretKeySecure = isFlagPassed("us")
	c.wasSetDataSecretKeySecure = isFlagPassed("ds")

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
