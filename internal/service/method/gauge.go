package method

import (
	"context"
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type Gauge struct {
	storage interfaces.Storage
	logger  interfaces.Logger
}

func NewGauge(storage interfaces.Storage, logger interfaces.Logger) *Gauge {
	return &Gauge{
		storage: storage,
		logger:  logger,
	}
}

func (g *Gauge) Update(ctx context.Context, updateParams params.Update) error {
	if err := g.storage.SetGauge(ctx, updateParams.ID, *updateParams.Value); err != nil {
		return fmt.Errorf("set gauge value: %w", err)
	}

	return nil
}

func (g *Gauge) Get(ctx context.Context, name string) (*models.Metrics, error) {
	gaugeValue, err := g.storage.GetGauge(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get gauge value: %w", err)
	}

	return &models.Metrics{
		ID:    name,
		MType: models.Gauge,
		Value: &gaugeValue,
	}, nil
}

func (g *Gauge) GetAll(ctx context.Context) (map[string]string, error) {
	gaugeValues, err := g.storage.GetAllGaugeValues(ctx)
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
