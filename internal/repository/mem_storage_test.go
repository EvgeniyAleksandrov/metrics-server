package repository_test

import (
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_CounterAddAnyElementsWithOneName_GetThisElementsAsSlice(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.CounterAdd("A", 1)
	require.NoError(t, err)

	err = sut.CounterAdd("A", 2)
	require.NoError(t, err)

	err = sut.CounterAdd("A", 3)
	require.NoError(t, err)

	values, err := sut.CounterGet("A")
	require.NoError(t, err)
	require.Equal(t, types.Counter(6), values)
}

func TestMemStorage_CounterAddAnyElementsWithIdenticalName_GetAllCountersByName(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.CounterAdd("A", 1)
	require.NoError(t, err)

	err = sut.CounterAdd("A", 2)
	require.NoError(t, err)

	err = sut.CounterAdd("B", 3)
	require.NoError(t, err)

	values, err := sut.CounterGet("A")
	require.NoError(t, err)
	require.Equal(t, types.Counter(3), values)

	values, err = sut.CounterGet("B")
	require.NoError(t, err)
	require.Equal(t, types.Counter(3), values)
}

func TestMemStorage_CounterGetUnavailableElement_ElementNotFoundError(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.CounterAdd("A", 1)
	require.NoError(t, err)

	err = sut.CounterAdd("A", 2)
	require.NoError(t, err)

	values, err := sut.CounterGet("B")
	require.ErrorIs(t, err, repository.ErrNotFoundElement)
	require.Empty(t, values)
}

func TestMemStorage_GaugeAddAnyElementsWithOneName_ElementWasUpdated(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.GaugeSet("A", 1.2)
	require.NoError(t, err)

	value, err := sut.GaugeGet("A")
	require.NoError(t, err)
	require.Equal(t, types.Gauge(1.2), value)

	err = sut.GaugeSet("A", 1.5)
	require.NoError(t, err)

	value, err = sut.GaugeGet("A")
	require.NoError(t, err)
	require.Equal(t, types.Gauge(1.5), value)
}

func TestMemStorage_AddGaugeElementGetCounterElementWithSameName_ElementWasNotFound(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.GaugeSet("A", 1.2)
	require.NoError(t, err)

	value, err := sut.CounterGet("A")
	require.ErrorIs(t, err, repository.ErrNotFoundElement)
	require.Empty(t, value)
}

func TestMemStorage_AddCounterElementGetGaugeElementWithSameName_ElementWasNotFound(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.CounterAdd("A", 1)
	require.NoError(t, err)

	value, err := sut.GaugeGet("A")
	require.ErrorIs(t, err, repository.ErrNotFoundElement)
	require.Empty(t, value)
}
