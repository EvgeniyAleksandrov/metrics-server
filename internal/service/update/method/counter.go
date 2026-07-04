//go:generate mockgen -source=counter.go -destination=counter_mock_test.go -package=method_test
package method

import (
	"fmt"
	"log"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/update"
)

type CounterStorage interface {
	CounterAdd(name string, value types.Counter) error
}

type Counter struct {
	storage CounterStorage
}

func NewCounter(storage CounterStorage) *Counter {
	return &Counter{
		storage: storage,
	}
}

func (c *Counter) Update(name, value string) error {
	preparedValue, err := types.CounterFromString(value)
	if err != nil {
		log.Printf("Counter converter failed: %s", err.Error())
		return update.ErrInvalidValueFormat
	}

	if err := c.storage.CounterAdd(name, preparedValue); err != nil {
		return fmt.Errorf("add value: %w", err)
	}

	return nil
}

func (c *Counter) GetType() metric.ValueType {
	return metric.Counter
}
