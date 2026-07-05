package repository

import (
	"errors"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
)

var ErrNotFoundElement = errors.New("not found element")

type MemStorage struct {
	gauge    map[string]types.Gauge
	counters map[string]types.Counter
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:    make(map[string]types.Gauge),
		counters: make(map[string]types.Counter),
	}
}

func (s *MemStorage) GaugeSet(name string, value types.Gauge) error {
	s.gauge[name] = value
	return nil
}

func (s *MemStorage) CounterAdd(name string, value types.Counter) error {
	if _, ok := s.counters[name]; !ok {
		s.counters[name] = 0
	}

	s.counters[name] += value

	return nil
}

func (s *MemStorage) GaugeGet(name string) (types.Gauge, error) {
	if value, ok := s.gauge[name]; ok {
		return value, nil
	}

	return 0, ErrNotFoundElement
}

func (s *MemStorage) CounterGet(name string) (types.Counter, error) {
	if values, ok := s.counters[name]; ok {
		return values, nil
	}

	return 0, ErrNotFoundElement
}

func (s *MemStorage) GetAllGaugeValues() (map[string]types.Gauge, error) {
	return s.gauge, nil
}

func (s *MemStorage) GetAllCounterValues() (map[string]types.Counter, error) {
	return s.counters, nil
}
