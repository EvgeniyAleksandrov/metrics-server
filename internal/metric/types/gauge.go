package types

import (
	"fmt"
	"strconv"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
)

type Gauge float64

func GaugeFromString(s string) (Gauge, error) {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parse string to float64: %w", err)
	}

	return Gauge(v), nil
}

func (g Gauge) String() string {
	return strconv.FormatFloat(float64(g), 'f', -1, 64)
}

func (g Gauge) Type() metric.ValueType {
	return metric.Gauge
}
