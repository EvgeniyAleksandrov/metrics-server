package repository_test

import (
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/stretchr/testify/require"
)

func TestMemStorage_CounterAddAnyElementsWithOneName_GetThisElementsAsSlice(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.AddCounter("A", 1)
	require.NoError(t, err)

	err = sut.AddCounter("A", 2)
	require.NoError(t, err)

	err = sut.AddCounter("A", 3)
	require.NoError(t, err)

	values, err := sut.GetCounter("A")
	require.NoError(t, err)
	require.Equal(t, int64(6), values)
}

func TestMemStorage_CounterAddAnyElementsWithIdenticalName_GetAllCountersByName(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.AddCounter("A", 1)
	require.NoError(t, err)

	err = sut.AddCounter("A", 2)
	require.NoError(t, err)

	err = sut.AddCounter("B", 3)
	require.NoError(t, err)

	values, err := sut.GetCounter("A")
	require.NoError(t, err)
	require.Equal(t, int64(3), values)

	values, err = sut.GetCounter("B")
	require.NoError(t, err)
	require.Equal(t, int64(3), values)
}

func TestMemStorage_CounterGetUnavailableElement_ElementNotFoundError(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.AddCounter("A", 1)
	require.NoError(t, err)

	err = sut.AddCounter("A", 2)
	require.NoError(t, err)

	values, err := sut.GetCounter("B")
	require.ErrorIs(t, err, repository.ErrNotFoundElement)
	require.Empty(t, values)
}

func TestMemStorage_GaugeAddAnyElementsWithOneName_ElementWasUpdated(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.SetGauge("A", 1.2)
	require.NoError(t, err)

	value, err := sut.GetGauge("A")
	require.NoError(t, err)
	require.Equal(t, 1.2, value)

	err = sut.SetGauge("A", 1.5)
	require.NoError(t, err)

	value, err = sut.GetGauge("A")
	require.NoError(t, err)
	require.Equal(t, 1.5, value)
}

func TestMemStorage_AddGaugeElementGetCounterElementWithSameName_ElementWasNotFound(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.SetGauge("A", 1.2)
	require.NoError(t, err)

	value, err := sut.GetCounter("A")
	require.ErrorIs(t, err, repository.ErrNotFoundElement)
	require.Empty(t, value)
}

func TestMemStorage_AddCounterElementGetGaugeElementWithSameName_ElementWasNotFound(t *testing.T) {
	sut := repository.NewMemStorage()

	err := sut.AddCounter("A", 1)
	require.NoError(t, err)

	value, err := sut.GetGauge("A")
	require.ErrorIs(t, err, repository.ErrNotFoundElement)
	require.Empty(t, value)
}
