package getter_test

import (
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/getter"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
	"github.com/stretchr/testify/require"
)

func TestPullCounter_GetAllWithCorrectResource_ReturnMetrics(t *testing.T) {
	sut := getter.NewPullCounter()

	metrics, err := sut.GetAll(resource.NewPullCounter())
	require.NoError(t, err)
	require.Len(t, metrics, 1)

	testedMetic := metrics[0]
	require.Equal(t, metric.PullCount, testedMetic.Name())
	require.Equal(t, types.Counter(0), testedMetic.Value())
}

func TestPullCounter_GetAllIncorrectResource_ReturnIncorrectResourceError(t *testing.T) {
	sut := getter.NewMemory()

	metrics, err := sut.GetAll(resource.NewRandom())
	require.ErrorIs(t, getter.ErrInvalidResource, err)
	require.Nil(t, metrics)
}
