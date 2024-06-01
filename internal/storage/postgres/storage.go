package postgres

import (
	"context"
	"fmt"

	"github.com/KartoonYoko/gophkeeper/internal/common"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	pool *pgxpool.Pool

	conf             Config
	secretkeyHandler *common.SecretKeyHandler
}

func New(ctx context.Context, conf Config) (*Storage, error) {
	s := new(Storage)

	s.conf = conf

	pool, err := pgxpool.New(context.Background(), s.conf.ConnectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}
	s.pool = pool

	db := stdlib.OpenDBFromPool(pool)
	err = migrate(db)
	if err != nil {
		return nil, fmt.Errorf("unable migrate database: %v", err)
	}

	s.secretkeyHandler, err = common.NewSecretKeyHandler(s.conf.SecretKeySecure)
	if err != nil {
		return nil, fmt.Errorf("unable to create secret key handler: %v", err)
	}

	return s, nil
}
