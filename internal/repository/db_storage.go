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

type Retryer interface {
	Do(ctx context.Context, fn func() error) error
}

type DBStorage struct {
	db      *sql.DB
	logger  interfaces.Logger
	retrier Retryer
}

func NewDBStorage(db *sql.DB, logger interfaces.Logger, retryer Retryer) *DBStorage {
	return &DBStorage{
		db:      db,
		logger:  logger,
		retrier: retryer,
	}
}

func (s *DBStorage) SetGauge(ctx context.Context, gauge values.Gauge) error {
	return s.retrier.Do(ctx, func() error {
		return s.setGauge(ctx, gauge)
	})
}

func (s *DBStorage) setGauge(ctx context.Context, gauge values.Gauge) error {
	_, err := s.db.ExecContext(ctx, insertGaugeSQL, gauge.Name, gauge.Value)
	if err != nil {
		return fmt.Errorf("put gauge value to db: %w", err)
	}

	return nil
}

func (s *DBStorage) SetGauges(ctx context.Context, gauges []values.Gauge) error {
	return s.retrier.Do(ctx, func() error {
		return s.setGauges(ctx, gauges)
	})
}

func (s *DBStorage) setGauges(ctx context.Context, gauges []values.Gauge) error {
	if len(gauges) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	defer tx.Rollback()

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
	return s.retrier.Do(ctx, func() error {
		return s.addCounter(ctx, counter)
	})
}

func (s *DBStorage) addCounter(ctx context.Context, counter values.Counter) error {
	_, err := s.db.ExecContext(ctx, insertCounterSQL, counter.Name, counter.Delta)
	if err != nil {
		return fmt.Errorf("add counter to db: %w", err)
	}

	return nil
}

func (s *DBStorage) AddCounters(ctx context.Context, counters []values.Counter) error {
	return s.retrier.Do(ctx, func() error {
		return s.addCounters(ctx, counters)
	})
}

func (s *DBStorage) addCounters(ctx context.Context, counters []values.Counter) error {
	if len(counters) == 0 {
		return nil
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create transaction: %w", err)
	}

	defer tx.Rollback()

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

	return value, s.retrier.Do(ctx, func() error {
		v, err := s.getGauge(ctx, name)

		value = v

		return err
	})
}

func (s *DBStorage) getGauge(ctx context.Context, name string) (float64, error) {
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

	return delta, s.retrier.Do(ctx, func() error {
		d, err := s.getCounter(ctx, name)

		delta = d

		return err
	})
}

func (s *DBStorage) getCounter(ctx context.Context, name string) (int64, error) {
	var delta int64

	if err := s.db.QueryRowContext(ctx, selectCounterSQL, name).Scan(&delta); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, ErrNotFoundElement
		}

		return 0, fmt.Errorf("select counter value: %w", err)
	}

	return delta, nil
}

func (s *DBStorage) GetAllGaugeValues(ctx context.Context) (map[string]float64, error) {
	var resultMap map[string]float64

	return resultMap, s.retrier.Do(ctx, func() error {
		m, err := s.getAllGaugeValues(ctx)
		resultMap = m

		return err
	})
}

func (s *DBStorage) getAllGaugeValues(ctx context.Context) (map[string]float64, error) {
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
	var resultMap map[string]int64

	return resultMap, s.retrier.Do(ctx, func() error {
		m, err := s.getAllCounterValues(ctx)
		resultMap = m

		return err
	})
}

func (s *DBStorage) getAllCounterValues(ctx context.Context) (map[string]int64, error) {
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
