package db

import (
	"database/sql"
	"log/slog"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/kylecain/wheel-of-wonder/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase(config *config.Config) *sql.DB {
	m, err := migrate.New(config.MigrationUrl, config.DatabaseUrl)
	if err != nil {
		slog.Error("migration setup failed", "error", err)
		os.Exit(1)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		slog.Error("migration failed", "error", err)
		os.Exit(1)
	}

	slog.Info("migrations applied successfully")

	dbPath := strings.TrimPrefix(config.DatabaseUrl, "sqlite3://")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		slog.Error("error creating database", "error", err)
		os.Exit(1)
	}

	return db
}
