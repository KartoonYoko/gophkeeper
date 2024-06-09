package clientstorage

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type migrationLogger struct{}

func (l *migrationLogger) Fatalf(format string, v ...interface{}) {
	fmt.Printf(format, v...)
}

func (l *migrationLogger) Printf(format string, v ...interface{}) {
	// не логируем выполнение миграции
}

func migrate(db *sql.DB) error {
	goose.SetBaseFS(embedMigrations)
	goose.SetLogger(new(migrationLogger))

	if err := goose.SetDialect("sqlite3"); err != nil {
		return fmt.Errorf("can not set dialect sqlite3: %w", err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		return err
	}

	return nil
}
