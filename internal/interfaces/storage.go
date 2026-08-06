//go:generate mockgen -source=storage.go -destination=storage_mock_test.go -package=interfaces_test
package interfaces

import (
	"context"
)

type Storage interface {
	SetGauge(ctx context.Context, name string, value float64) error
	AddCounter(ctx context.Context, name string, value int64) error
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
