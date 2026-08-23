package getter

import (
	"fmt"
	"log"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
)

type System struct{}

func NewSystem() *System {
	return &System{}
}

func (s *System) Resource() resource.Type {
	return resource.SystemType
}

func (s *System) GetAll(r resource.Resource) ([]*metric.Metric, error) {
	systemStats, ok := r.Get().(resource.SystemStats)
	if !ok {
		log.Println(r.Type())
		return nil, ErrInvalidResource
	}

	cpuMetrics := s.createCpuMetrics(systemStats)

	memMetrics := []*metric.Metric{
		metric.NewMetric(metric.TotalMemory, types.Gauge(systemStats.MemTotal)),
		metric.NewMetric(metric.FreeMemory, types.Gauge(systemStats.FreeMemory)),
	}

	return append(memMetrics, cpuMetrics...), nil
}

func (s *System) createCpuMetrics(systemStats resource.SystemStats) []*metric.Metric {
	cpuUsageMetrics := make([]*metric.Metric, len(systemStats.CPUUsage))

	for i, value := range systemStats.CPUUsage {
		cpuUsageMetrics[i] = metric.NewMetric(
			metric.Name(fmt.Sprintf("CPUutilization%d", i+1)),
			types.Gauge(value),
		)
	}

	return cpuUsageMetrics
}
