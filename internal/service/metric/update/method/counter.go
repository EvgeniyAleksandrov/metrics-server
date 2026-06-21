package method

import (
	"fmt"
	"log"
	"strconv"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/metric"
)

type CounterStorage interface {
	CounterAdd(name string, value int64) error
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
	preparedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		log.Printf("Int64 converter failed: %s", err.Error())
		return metric.ErrInvalidValueFormat
	}

	if err := c.storage.CounterAdd(name, preparedValue); err != nil {
		return fmt.Errorf("add value: %w", err)
	}

	return nil
}

func (c *Counter) GetType() string {
	return "counter"
}
