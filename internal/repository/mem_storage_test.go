package repository_test

import (
	"context"
	"testing"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/stretchr/testify/suite"
)

type MemStorageSuite struct {
	suite.Suite

	ctx context.Context
}

func TestMemStorageSuite(t *testing.T) {
	t.Parallel()

	suite.Run(t, &MemStorageSuite{ctx: context.Background()})
}

func (s *MemStorageSuite) TestMemStorage_CounterAddAnyElementsWithOneName_GetThisElementsAsSlice() {
	sut := repository.NewMemStorage()

	err := sut.AddCounter(s.ctx, "A", 1)
	s.Require().NoError(err)

	err = sut.AddCounter(s.ctx, "A", 2)
	s.Require().NoError(err)

	err = sut.AddCounter(s.ctx, "A", 3)
	s.Require().NoError(err)

	values, err := sut.GetCounter(s.ctx, "A")
	s.Require().NoError(err)
	s.Require().Equal(int64(6), values)
}

func (s *MemStorageSuite) TestMemStorage_CounterAddAnyElementsWithIdenticalName_GetAllCountersByName() {
	sut := repository.NewMemStorage()

	err := sut.AddCounter(s.ctx, "A", 1)
	s.Require().NoError(err)

	err = sut.AddCounter(s.ctx, "A", 2)
	s.Require().NoError(err)

	err = sut.AddCounter(s.ctx, "B", 3)
	s.Require().NoError(err)

	values, err := sut.GetCounter(s.ctx, "A")
	s.Require().NoError(err)
	s.Require().Equal(int64(3), values)

	values, err = sut.GetCounter(s.ctx, "B")
	s.Require().NoError(err)
	s.Require().Equal(int64(3), values)
}

func (s *MemStorageSuite) TestMemStorage_CounterGetUnavailableElement_ElementNotFoundError() {
	sut := repository.NewMemStorage()

	err := sut.AddCounter(s.ctx, "A", 1)
	s.Require().NoError(err)

	err = sut.AddCounter(s.ctx, "A", 2)
	s.Require().NoError(err)

	values, err := sut.GetCounter(s.ctx, "B")
	s.Require().ErrorIs(err, repository.ErrNotFoundElement)
	s.Require().Empty(values)
}

func (s *MemStorageSuite) TestMemStorage_GaugeAddAnyElementsWithOneName_ElementWasUpdated() {
	sut := repository.NewMemStorage()

	err := sut.SetGauge(s.ctx, "A", 1.2)
	s.Require().NoError(err)

	value, err := sut.GetGauge(s.ctx, "A")
	s.Require().NoError(err)
	s.Require().Equal(1.2, value)

	err = sut.SetGauge(s.ctx, "A", 1.5)
	s.Require().NoError(err)

	value, err = sut.GetGauge(s.ctx, "A")
	s.Require().NoError(err)
	s.Require().Equal(1.5, value)
}

func (s *MemStorageSuite) TestMemStorage_AddGaugeElementGetCounterElementWithSameName_ElementWasNotFound() {
	sut := repository.NewMemStorage()

	err := sut.SetGauge(s.ctx, "A", 1.2)
	s.Require().NoError(err)

	value, err := sut.GetCounter(s.ctx, "A")
	s.Require().ErrorIs(err, repository.ErrNotFoundElement)
	s.Require().Empty(value)
}

func (s *MemStorageSuite) TestMemStorage_AddCounterElementGetGaugeElementWithSameName_ElementWasNotFound() {
	sut := repository.NewMemStorage()

	err := sut.AddCounter(s.ctx, "A", 1)
	s.Require().NoError(err)

	value, err := sut.GetGauge(s.ctx, "A")
	s.Require().ErrorIs(err, repository.ErrNotFoundElement)
	s.Require().Empty(value)
}
