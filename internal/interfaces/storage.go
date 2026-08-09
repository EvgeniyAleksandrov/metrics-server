//go:generate mockgen -source=storage.go -destination=storage_mock_test.go -package=interfaces_test
package interfaces

import (
	"context"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/repository/values"
)

type Storage interface {
	SetGauge(ctx context.Context, gauge values.Gauge) error
	SetGauges(ctx context.Context, gauges []values.Gauge) error
	AddCounters(ctx context.Context, counter []values.Counter) error
	AddCounter(ctx context.Context, counter values.Counter) error
	Ping(ctx context.Context) error
	GetGauge(ctx context.Context, name string) (float64, error)
	GetCounter(ctx context.Context, name string) (int64, error)
	GetAllGaugeValues(ctx context.Context) (map[string]float64, error)
	GetAllCounterValues(ctx context.Context) (map[string]int64, error)
	Close()
}

type MarshaledStorage interface {
	Storage
	Marshal() ([]byte, error)
	Unmarshal(jsonData []byte) error
}
