package getter

import (
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
)

type Random struct{}

func NewRandom() *Random {
	return &Random{}
}

func (r *Random) Resource() resource.Type {
	return resource.RandomType
}

func (r *Random) GetAll(res resource.Resource) ([]*metric.Metric, error) {
	randomValue, ok := res.Get().(float64)
	if !ok {
		return nil, ErrInvalidResource
	}

	return []*metric.Metric{
		metric.NewMetric(metric.RandomValue, types.Gauge(randomValue)),
	}, nil
}
