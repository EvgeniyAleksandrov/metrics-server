package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

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

func (s *MemStorage) SetGauge(_ context.Context, name string, value float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.setGauge(name, value)
}

func (s *MemStorage) setGauge(name string, value float64) error {
	s.gauge[name] = value
	return nil
}

func (s *MemStorage) AddCounter(_ context.Context, name string, value int64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.addCounter(name, value)
}

func (s *MemStorage) addCounter(name string, value int64) error {
	if _, ok := s.counters[name]; !ok {
		s.counters[name] = value
		return nil
	}

	s.counters[name] += value

	return nil
}

func (s *MemStorage) Ping(_ context.Context) error {
	return nil
}

func (s *MemStorage) GetGauge(_ context.Context, name string) (float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if value, ok := s.gauge[name]; ok {
		return value, nil
	}

	return 0, ErrNotFoundElement
}

func (s *MemStorage) GetCounter(_ context.Context, name string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if values, ok := s.counters[name]; ok {
		return values, nil
	}

	return 0, ErrNotFoundElement
}

func (s *MemStorage) GetAllGaugeValues(_ context.Context) (map[string]float64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.gauge, nil
}

func (s *MemStorage) GetAllCounterValues(_ context.Context) (map[string]int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	return s.counters, nil
}

func (s *MemStorage) Marshal() ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	metricsSlice := make([]models.Metrics, 0, len(s.gauge)+len(s.counters))

	for name, gValue := range s.gauge {
		metricsSlice = append(metricsSlice,
			models.Metrics{
				ID:    name,
				MType: models.Gauge,
				Value: &gValue,
			},
		)
	}

	for name, cValue := range s.counters {
		metricsSlice = append(metricsSlice,
			models.Metrics{
				ID:    name,
				MType: models.Counter,
				Delta: &cValue,
			},
		)
	}

	jsonData, err := json.Marshal(metricsSlice)
	if err != nil {
		return nil, fmt.Errorf("marshal metrics: %w", err)
	}

	return jsonData, nil
}

func (s *MemStorage) Unmarshal(jsonData []byte) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var metricSlice []models.Metrics

	if err := json.Unmarshal(jsonData, &metricSlice); err != nil {
		return fmt.Errorf("unmarshal json data: %w", err)
	}

	for _, metric := range metricSlice {
		switch metric.MType {
		case models.Gauge:
			if metric.Value == nil {
				continue
			}

			_ = s.setGauge(metric.ID, *metric.Value)
		case models.Counter:
			if metric.Delta == nil {
				continue
			}

			_ = s.addCounter(metric.ID, *metric.Delta)
		}
	}

	return nil
}

func (s *MemStorage) Close() {}
