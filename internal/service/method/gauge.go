//go:generate mockgen -source=gauge.go -destination=gauge_mock_test.go -package=method_test
package method

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type GaugeStorage interface {
	SetGauge(name string, value float64) error
	GetGauge(name string) (float64, error)
	GetAllGaugeValues() (map[string]float64, error)
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

func (g *Gauge) Update(updateParams params.Update) error {
	if err := g.storage.SetGauge(updateParams.ID, *updateParams.Value); err != nil {
		return fmt.Errorf("set gauge value: %w", err)
	}

	return nil
}

func (g *Gauge) Get(name string) (*models.Metrics, error) {
	gaugeValue, err := g.storage.GetGauge(name)
	if err != nil {
		return nil, fmt.Errorf("get gauge value: %w", err)
	}

	return &models.Metrics{
		ID:    name,
		MType: models.Gauge,
		Value: &gaugeValue,
	}, nil
}

func (g *Gauge) GetAll() (map[string]string, error) {
	gaugeValues, err := g.storage.GetAllGaugeValues()
	if err != nil {
		return nil, fmt.Errorf("get gauge value: %w", err)
	}

	values := make(map[string]string, len(gaugeValues))

	for name, value := range gaugeValues {
		values[name] = types.Gauge(value).String()
	}

	return values, nil
}

func (g *Gauge) GetType() string {
	return models.Gauge
}
