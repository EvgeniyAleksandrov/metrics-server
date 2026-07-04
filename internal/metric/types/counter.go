package types

import (
	"fmt"
	"strconv"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
)

type Counter int64

func CounterFromString(s string) (Counter, error) {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse string to int64: %w", err)
	}

	return Counter(v), nil
}

func (c Counter) String() string {
	return fmt.Sprintf("%d", c)
}

func (c Counter) Type() metric.ValueType {
	return metric.Counter
}
