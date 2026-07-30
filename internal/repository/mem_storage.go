package repository

import (
	"errors"
	"sync"
)

var ErrNotFoundElement = errors.New("not found element")

type MemStorage struct {
	gauge    map[string]float64
	counters map[string]int64

	mu sync.RWMutex
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:    make(map[string]float64),
		counters: make(map[string]int64),
	}
}

func (s *MemStorage) SetGauge(name string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.gauge[name] = value
	return nil
}

func (s *MemStorage) AddCounter(name string, value int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.counters[name]; !ok {
		s.counters[name] = value
		return nil
	}

	s.counters[name] += value

	return nil
}

func (s *MemStorage) GetGauge(name string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if value, ok := s.gauge[name]; ok {
		return value, nil
	}

	return 0, ErrNotFoundElement
}

func (s *MemStorage) GetCounter(name string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if values, ok := s.counters[name]; ok {
		return values, nil
	}

	return 0, ErrNotFoundElement
}

func (s *MemStorage) GetAllGaugeValues() (map[string]float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.gauge, nil
}

func (s *MemStorage) GetAllCounterValues() (map[string]int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.counters, nil
}
