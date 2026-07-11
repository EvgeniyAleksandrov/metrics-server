package getter_test

import (
	"slices"
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/getter"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
	"github.com/stretchr/testify/require"
)

func TestRandom_GetAllWithCorrectResource_ReturnMetrics(t *testing.T) {
	var expectedMetricsName = []metric.Name{
		metric.RandomValue,
	}

	sut := getter.NewRandom()

	metrics, err := sut.GetAll(resource.NewRandom())
	require.NoError(t, err)
	for _, currentMetric := range metrics {
		require.True(t, slices.Contains(expectedMetricsName, currentMetric.Name()))
	}
}

func TestRandom_GetAllIncorrectResource_ReturnIncorrectResourceError(t *testing.T) {
	sut := getter.NewRandom()

	metrics, err := sut.GetAll(resource.NewMemory())
	require.ErrorIs(t, getter.ErrInvalidResource, err)
	require.Nil(t, metrics)
}
