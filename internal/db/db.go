package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/kylecain/wheel-of-wonder/internal/config"
	_ "github.com/mattn/go-sqlite3"
)

var (
	migrationUrl = "file://internal/db/migrations"
	databaseUrl  = "sqlite3://data/db.sqlite3"
)

func NewDatabase(coneefig *config.Config) (*sql.DB, error) {
	m, err := migrate.New(migrationUrl, databaseUrl)
	if err != nil {
		return nil, fmt.Errorf("migration setup: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return nil, fmt.Errorf("applying migrations: %w", err)
	}

	dbPath := strings.TrimPrefix(databaseUrl, "sqlite3://")
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("open sqlite db %s: %w", dbPath, err)
	}

	return db, nil
}
