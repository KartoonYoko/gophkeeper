package postgres

import (
	"context"
	"fmt"

	"github.com/KartoonYoko/gophkeeper/internal/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

type Storage struct {
	pool *pgxpool.Pool

	conf Config
}

func New(ctx context.Context, conf Config) (*Storage, error) {
	s := new(Storage)

	s.conf = conf

	logger.Log.Info("trying connect to databse with dsn", zap.String("dsn", conf.ConnectionString))
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

	return s, nil
}
