//go:generate mockgen -source=gauge.go -destination=gauge_mock_test.go -package=method_test
package method

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
	"go.uber.org/zap"
)

type GaugeStorage interface {
	GaugeSet(name string, value types.Gauge) error
	GaugeGet(name string) (types.Gauge, error)
	GetAllGaugeValues() (map[string]types.Gauge, error)
}

type Gauge struct {
	storage GaugeStorage
	logger  Logger
}

func NewGauge(storage GaugeStorage, logger Logger) *Gauge {
	return &Gauge{
		storage: storage,
		logger:  logger,
	}
}

func (g *Gauge) Update(name, value string) error {
	preparedValue, err := types.GaugeFromString(value)
	if err != nil {
		g.logger.Warn("Gauge converter failed", zap.Error(err))
		return service.ErrInvalidValueFormat
	}

	if err := g.storage.GaugeSet(name, preparedValue); err != nil {
		return fmt.Errorf("set value: %w", err)
	}

	return nil
}

func (g *Gauge) Get(name string) (string, error) {
	value, err := g.storage.GaugeGet(name)
	if err != nil {
		return "", fmt.Errorf("get gauge value: %w", err)
	}

	return value.String(), nil
}

func (g *Gauge) GetAll() (map[string]string, error) {
	gaugeValues, err := g.storage.GetAllGaugeValues()
	if err != nil {
		return nil, fmt.Errorf("get gauge value: %w", err)
	}

	values := make(map[string]string, len(gaugeValues))

	for name, value := range gaugeValues {
		values[name] = value.String()
	}

	return values, nil
}

func (g *Gauge) GetType() metric.ValueType {
	return metric.Gauge
}
