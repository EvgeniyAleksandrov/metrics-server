package repository

type MemStorage struct {
	gauge   map[string]float64
	counter map[string][]int64
}

func NewMemStorage() *MemStorage {
	return &MemStorage{
		gauge:   make(map[string]float64),
		counter: make(map[string][]int64),
	}
}

func (s *MemStorage) GaugeSet(name string, value float64) error {
	s.gauge[name] = value
	return nil
}

func (s *MemStorage) CounterAdd(name string, value int64) error {
	if _, ok := s.counter[name]; !ok {
		s.counter[name] = make([]int64, 1)
	}

	s.counter[name] = append(s.counter[name], value)

	return nil
}
