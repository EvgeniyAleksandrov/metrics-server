package translation

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type Counter struct{}

func NewCounter() *Counter {
	return &Counter{}
}

func (g *Counter) Type() metric.ValueType {
	return metric.Counter
}

func (g *Counter) Translate(m *metric.Metric) (*models.Metrics, error) {
	if m == nil {
		return nil, ErrEmptyMetric
	}

	cValue, err := types.CounterFromString(fmt.Sprintf("%s", m.Value()))
	if err != nil {
		return nil, fmt.Errorf("translate gauge: %w", err)
	}

	value := int64(cValue)
	return &models.Metrics{
		ID:    string(m.Name()),
		MType: models.Counter,
		Delta: &value,
	}, nil
}
