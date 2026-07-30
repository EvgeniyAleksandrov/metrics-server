package translation

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type Gauge struct{}

func NewGauge() *Gauge {
	return &Gauge{}
}

func (g *Gauge) Type() metric.ValueType {
	return metric.Gauge
}

func (g *Gauge) Translate(m *metric.Metric) (*models.Metrics, error) {
	if m == nil {
		return nil, ErrEmptyMetric
	}

	gValue, err := types.GaugeFromString(fmt.Sprintf("%s", m.Value()))
	if err != nil {
		return nil, fmt.Errorf("translate gauge: %w", err)
	}

	value := float64(gValue)
	return &models.Metrics{
		ID:    string(m.Name()),
		MType: models.Gauge,
		Value: &value,
	}, nil
}
