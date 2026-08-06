package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var createTableSQL = `
CREATE TABLE IF NOT EXISTS metrics (
    id    TEXT PRIMARY KEY,
    type  TEXT NOT NULL,
    delta BIGINT,
    value DOUBLE PRECISION,
    hash  TEXT,

    CONSTRAINT metrics_type_value_check
        CHECK (
            (type = 'counter' AND delta IS NOT NULL AND value IS NULL)
            OR
            (type = 'gauge'   AND value IS NOT NULL AND delta IS NULL)
        )
);
`

func BuildDBConnection(ctx context.Context, dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("open database: %w", err)
	}

	pingContext, cncl := context.WithTimeout(ctx, 1*time.Second)
	defer cncl()

	if err = db.PingContext(pingContext); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	createContext, cancel := context.WithTimeout(ctx, time.Duration(5)*time.Second)
	defer cancel()

	if _, err := db.ExecContext(createContext, createTableSQL); err != nil {
		return nil, fmt.Errorf("create tables: %w", err)
	}

	return db, nil
}
