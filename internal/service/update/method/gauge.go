//go:generate mockgen -source=gauge.go -destination=gauge_mock_test.go -package=method_test
package method

import (
	"fmt"
	"log"

	metric2 "github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/update"
)

type GaugeStorage interface {
	GaugeSet(name string, value types.Gauge) error
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
	preparedValue, err := types.GaugeFromString(value)
	if err != nil {
		log.Printf("Float64 converter failed: %s", err.Error())
		return update.ErrInvalidValueFormat
	}

	if err := g.storage.GaugeSet(name, preparedValue); err != nil {
		return fmt.Errorf("set value: %w", err)
	}

	return nil
}

func (g *Gauge) GetType() metric2.ValueType {
	return metric2.Gauge
}
