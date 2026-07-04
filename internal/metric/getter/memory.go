package getter

import (
	"runtime"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
)

type Memory struct{}

func NewMemory() *Memory {
	return &Memory{}
}

func (a *Memory) Resource() resource.Type {
	return resource.MemoryType
}

func (a *Memory) GetAll(r resource.Resource) ([]*metric.Metric, error) {
	memStats, ok := r.Get().(runtime.MemStats)
	if !ok {
		return nil, ErrInvalidResource
	}

	return []*metric.Metric{
		metric.NewMetric(metric.Alloc, types.Gauge(memStats.Alloc)),
		metric.NewMetric(metric.BuckHashSys, types.Gauge(memStats.BuckHashSys)),
		metric.NewMetric(metric.Frees, types.Counter(memStats.Frees)),
		metric.NewMetric(metric.GCCPUFraction, types.Gauge(memStats.GCCPUFraction)),
		metric.NewMetric(metric.GCSys, types.Gauge(memStats.GCSys)),
		metric.NewMetric(metric.HeapAlloc, types.Gauge(memStats.HeapAlloc)),
		metric.NewMetric(metric.HeapIdle, types.Gauge(memStats.HeapIdle)),
		metric.NewMetric(metric.HeapInuse, types.Gauge(memStats.HeapInuse)),
		metric.NewMetric(metric.HeapObjects, types.Gauge(memStats.HeapObjects)),
		metric.NewMetric(metric.HeapReleased, types.Gauge(memStats.HeapReleased)),
		metric.NewMetric(metric.HeapSys, types.Gauge(memStats.HeapSys)),
		metric.NewMetric(metric.LastGC, types.Gauge(memStats.LastGC)),
		metric.NewMetric(metric.Lookups, types.Counter(memStats.Lookups)),
		metric.NewMetric(metric.MCacheInuse, types.Gauge(memStats.MCacheInuse)),
		metric.NewMetric(metric.MCacheSys, types.Gauge(memStats.MCacheSys)),
		metric.NewMetric(metric.MSpanInuse, types.Gauge(memStats.MSpanInuse)),
		metric.NewMetric(metric.MSpanSys, types.Gauge(memStats.MSpanSys)),
		metric.NewMetric(metric.Mallocs, types.Counter(memStats.Mallocs)),
		metric.NewMetric(metric.NextGC, types.Gauge(memStats.NextGC)),
		metric.NewMetric(metric.NumForcedGC, types.Counter(uint64(memStats.NumForcedGC))),
		metric.NewMetric(metric.NumGC, types.Counter(uint64(memStats.NumGC))),
		metric.NewMetric(metric.OtherSys, types.Gauge(memStats.OtherSys)),
		metric.NewMetric(metric.PauseTotalNs, types.Counter(memStats.PauseTotalNs)),
		metric.NewMetric(metric.StackInuse, types.Gauge(memStats.StackInuse)),
		metric.NewMetric(metric.StackSys, types.Gauge(memStats.StackSys)),
		metric.NewMetric(metric.Sys, types.Gauge(memStats.Sys)),
		metric.NewMetric(metric.TotalAlloc, types.Counter(memStats.TotalAlloc)),
	}, nil
}
