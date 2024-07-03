package postgres

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/KartoonYoko/gophkeeper/internal/logger"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type migrationLogger struct{}

func (l *migrationLogger) Fatalf(format string, v ...interface{}) {
	logger.Log.Sugar().Errorf(format, v)
}

func (l *migrationLogger) Printf(format string, v ...interface{}) {
	logger.Log.Sugar().Infof(format, v)
}

func migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(new(migrationLogger))

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("postgres migrate set dialect postgres: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
