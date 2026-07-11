package getter_test

import (
	"slices"
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/getter"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
	"github.com/stretchr/testify/require"
)

func TestMemory_GetAllWithCorrectResource_ReturnMetrics(t *testing.T) {
	var expectedMetricsName = []metric.Name{
		metric.Alloc,
		metric.BuckHashSys,
		metric.Frees,
		metric.GCCPUFraction,
		metric.GCSys,
		metric.HeapAlloc,
		metric.HeapIdle,
		metric.HeapInuse,
		metric.HeapObjects,
		metric.HeapReleased,
		metric.HeapSys,
		metric.LastGC,
		metric.Lookups,
		metric.MCacheInuse,
		metric.MCacheSys,
		metric.MSpanInuse,
		metric.MSpanSys,
		metric.Mallocs,
		metric.NextGC,
		metric.NumForcedGC,
		metric.NumGC,
		metric.OtherSys,
		metric.PauseTotalNs,
		metric.StackInuse,
		metric.StackSys,
		metric.Sys,
		metric.TotalAlloc,
	}

	sut := getter.NewMemory()

	metrics, err := sut.GetAll(resource.NewMemory())
	require.NoError(t, err)
	for _, currentMetric := range metrics {
		require.True(t, slices.Contains(expectedMetricsName, currentMetric.Name()))

	}
}

func TestMemory_GetAllIncorrectResource_ReturnIncorrectResourceError(t *testing.T) {
	sut := getter.NewMemory()

	metrics, err := sut.GetAll(resource.NewRandom())
	require.ErrorIs(t, getter.ErrInvalidResource, err)
	require.Nil(t, metrics)
}
