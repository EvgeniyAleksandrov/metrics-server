package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/migrations"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func BuildDBConnection(ctx context.Context, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	pingContext, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()

	if err = db.PingContext(pingContext); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	if err := migrateDB(db); err != nil {
		return nil, fmt.Errorf("migrate db: %w", err)
	}

	return db, nil
}

func migrateDB(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("create migration driver: %w", err)
	}

	sourceDrive, err := iofs.New(migrations.GetMigrationFS(), ".")
	if err != nil {
		return fmt.Errorf("create source drive: %w", err)
	}

	dbMigrate, err := migrate.NewWithInstance(
		"migration_source",
		sourceDrive,
		"pgx",
		driver,
	)
	if err != nil {
		return fmt.Errorf("create migrate: %w", err)
	}

	if err := dbMigrate.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("up migrations: %w", err)
	}

	return nil
}
