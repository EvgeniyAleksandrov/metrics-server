package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/repository/values"
	"go.uber.org/zap"
)

const (
	insertGaugeSQL = `
INSERT INTO metrics (id, type, value) 
VALUES ($1, 'gauge', $2)
ON CONFLICT (id)
DO UPDATE SET
value = EXCLUDED.value;
`
	insertCounterSQL = `
INSERT INTO metrics (id, type, delta) 
VALUES ($1, 'counter', $2)
ON CONFLICT (id)
DO UPDATE SET
delta = metrics.delta + EXCLUDED.delta;
`
	selectGaugeSQL = `
SELECT value FROM metrics
WHERE id = $1 AND type='gauge';
`
	selectCounterSQL = `
SELECT delta FROM metrics
WHERE id = $1 AND type='counter';
`

	selectAllGaugeSQL = `
SELECT id, value FROM metrics
WHERE type='gauge';
`
	selectAllCounterSQL = `
SELECT id, delta FROM metrics
WHERE type='counter';
`
)

type DBStorage struct {
	db *sql.DB

	logger interfaces.Logger
}

func NewDBStorage(db *sql.DB, logger interfaces.Logger) *DBStorage {
	return &DBStorage{
		db:     db,
		logger: logger,
	}
}

func (s *DBStorage) SetGauge(ctx context.Context, gauge values.Gauge) error {
	_, err := s.db.ExecContext(ctx, insertGaugeSQL, gauge.Name, gauge.Value)
	if err != nil {
		return fmt.Errorf("put gauge value to db: %w", err)
	}

	return nil
}

func (s *DBStorage) SetGauges(ctx context.Context, gauges []values.Gauge) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			s.logger.Warn("Rollback error", zap.Error(err))
		}
	}()

	stmt, err := tx.PrepareContext(ctx, insertGaugeSQL)
	if err != nil {
		return fmt.Errorf("prepare gauge insert sql: %w", err)
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			s.logger.Warn("Close Gauge stmt error", zap.Error(err))
		}
	}()

	for _, gauge := range gauges {
		if _, err := stmt.ExecContext(ctx, gauge.Name, gauge.Value); err != nil {
			return fmt.Errorf("put gauge value to db: %w", err)
		}
	}

	return tx.Commit()
}

func (s *DBStorage) AddCounter(ctx context.Context, counter values.Counter) error {
	_, err := s.db.ExecContext(ctx, insertCounterSQL, counter.Name, counter.Delta)
	if err != nil {
		return fmt.Errorf("add counter to db: %w", err)
	}

	return nil
}

func (s *DBStorage) AddCounters(ctx context.Context, counters []values.Counter) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	defer func() {
		if err := tx.Rollback(); err != nil {
			s.logger.Warn("Rollback error", zap.Error(err))
		}
	}()

	stmt, err := tx.PrepareContext(ctx, insertCounterSQL)
	if err != nil {
		return fmt.Errorf("prepare counter insert sql: %w", err)
	}

	defer func() {
		if err := stmt.Close(); err != nil {
			s.logger.Warn("Close Counter stmt error", zap.Error(err))
		}
	}()

	for _, gauge := range counters {
		if _, err := stmt.ExecContext(ctx, gauge.Name, gauge.Delta); err != nil {
			return fmt.Errorf("add counter value to db: %w", err)
		}
	}

	return tx.Commit()
}

func (s *DBStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	var value float64

	if err := s.db.QueryRowContext(ctx, selectGaugeSQL, name).Scan(&value); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFoundElement
		}

		return 0, fmt.Errorf("select gauge value: %w", err)
	}

	return value, nil
}

func (s *DBStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	var delta int64

	if err := s.db.QueryRowContext(ctx, selectCounterSQL, name).Scan(&delta); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFoundElement
		}

		return 0, fmt.Errorf("select cluters value: %w", err)
	}

	return delta, nil
}

func (s *DBStorage) GetAllGaugeValues(ctx context.Context) (map[string]float64, error) {
	rows, err := s.db.QueryContext(ctx, selectAllGaugeSQL)
	if err != nil {
		return nil, fmt.Errorf("select all gauge values: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			s.logger.Warn("Close rows error", zap.Error(err))
		}
	}()

	gauges := make(map[string]float64)

	var (
		id    string
		value float64
	)

	for rows.Next() {
		if err := rows.Scan(&id, &value); err != nil {
			return nil, fmt.Errorf("scan gauge values: %w", err)
		}

		gauges[id] = value
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get rows error: %w", err)
	}

	return gauges, nil
}

func (s *DBStorage) GetAllCounterValues(ctx context.Context) (map[string]int64, error) {
	rows, err := s.db.QueryContext(ctx, selectAllCounterSQL)
	if err != nil {
		return nil, fmt.Errorf("select all counter values: %w", err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			s.logger.Warn("Close rows error", zap.Error(err))
		}
	}()

	counters := make(map[string]int64)

	var (
		id    string
		delta int64
	)

	for rows.Next() {
		if err := rows.Scan(&id, &delta); err != nil {
			return nil, fmt.Errorf("scan counter values: %w", err)
		}

		counters[id] = delta
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get rows error: %w", err)
	}

	return counters, nil
}

func (s *DBStorage) Ping(ctx context.Context) error {
	if err := s.db.PingContext(ctx); err != nil {
		return fmt.Errorf("ping db: %w", err)
	}

	return nil
}

func (s *DBStorage) Close() {
	s.logger.Info("Close db connection")

	if err := s.db.Close(); err != nil {
		s.logger.Warn("Close db connection with error", zap.Error(err))
	}
}
