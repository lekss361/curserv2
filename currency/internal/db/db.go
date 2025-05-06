package db

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/lekss361/curserv2/currency/internal/config"
	_ "github.com/lib/pq"
	"log"
)

// New создаёт *sql.DB для PostgreSQL и проверяет соединение.
func New() (*sql.DB, error) {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("db.Ping: %w", err)
	}
	return db, nil
}

func Migrate(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres.WithInstance: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://C:/Users/User/Desktop/golangeduc/currencyservice/currency.go/internal/migrations",
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %w", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("m.Up: %w", err)
	}
	return nil
}

func MigrateDown(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("postgres.WithInstance: %w", err)
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file://C:/Users/User/Desktop/golangeduc/currencyservice/currency.go/internal/migrations",
		"postgres", driver,
	)
	if err != nil {
		return fmt.Errorf("migrate.NewWithDatabaseInstance: %w", err)
	}
	if err := m.Down(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("m.Down: %w", err)
	}
	return nil
}
