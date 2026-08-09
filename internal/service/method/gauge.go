package method

import (
	"context"
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/repository/values"
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
	if err := g.storage.SetGauge(ctx, values.Gauge{Name: updateParams.ID, Value: *updateParams.Value}); err != nil {
		return fmt.Errorf("set gauge value: %w", err)
	}

	return nil
}

func (g *Gauge) UpdateBatch(ctx context.Context, updates []params.Update) error {
	gauges := make([]values.Gauge, len(updates))

	for _, update := range updates {
		gauges = append(gauges, values.Gauge{Name: update.ID, Value: *update.Value})
	}

	if err := g.storage.SetGauges(ctx, gauges); err != nil {
		return fmt.Errorf("set gauges: %w", err)
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

	getValues := make(map[string]string, len(gaugeValues))

	for name, value := range gaugeValues {
		getValues[name] = types.Gauge(value).String()
	}

	return getValues, nil
}

func (g *Gauge) GetType() string {
	return models.Gauge
}
