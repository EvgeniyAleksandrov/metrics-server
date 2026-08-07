package getter

import (
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
)

type PullCounter struct{}

func NewPullCounter() *PullCounter {
	return &PullCounter{}
}

func (a *PullCounter) Resource() resource.Type {
	return resource.PullCounterType
}

func (a *PullCounter) GetAll(r resource.Resource) ([]*metric.Metric, error) {
	pCounter, ok := r.Get().(int64)
	if !ok {
		return nil, ErrInvalidResource
	}

	return []*metric.Metric{metric.NewMetric(metric.PollCount, types.Counter(pCounter))}, nil
}
