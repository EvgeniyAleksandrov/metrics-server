package method

import (
	"fmt"
	"log"
	"strconv"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/metric"
)

type GaugeStorage interface {
	GaugeSet(name string, value float64) error
}

type Gauge struct {
	storage GaugeStorage
}

func NewGauge(storage GaugeStorage) *Gauge {
	return &Gauge{
		storage: storage,
	}
}

func (g *Gauge) Update(name, value string) error {
	preparedValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		log.Printf("Float64 converter failed: %s", err.Error())
		return metric.ErrInvalidValueFormat
	}

	if err := g.storage.GaugeSet(name, preparedValue); err != nil {
		return fmt.Errorf("set value: %w", err)
	}

	return nil
}

func (g *Gauge) GetType() string {
	return "gauge"
}
