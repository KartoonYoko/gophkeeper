package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
)

type Storage struct {
	pool *pgxpool.Pool

	conf Config
}

func New(ctx context.Context, conf Config) (*Storage, error) {
	s := new(Storage)

	s.conf = conf

	pool, err := pgxpool.New(context.Background(), s.conf.connectionString)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %v", err)
	}

	db := stdlib.OpenDBFromPool(pool)
	err = migrate(db)
	if err != nil {
		return nil, fmt.Errorf("unable migrate database: %v", err)
	}

	return s, nil
}
