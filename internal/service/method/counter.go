//go:generate mockgen -source=counter.go -destination=counter_mock_test.go -package=method_test
package method

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
	"go.uber.org/zap"
)

type CounterStorage interface {
	CounterAdd(name string, value types.Counter) error
	CounterGet(name string) (types.Counter, error)
	GetAllCounterValues() (map[string]types.Counter, error)
}

type Counter struct {
	storage CounterStorage
	logger  Logger
}

func NewCounter(storage CounterStorage, logger Logger) *Counter {
	return &Counter{
		storage: storage,
		logger:  logger,
	}
}

func (c *Counter) Update(name, value string) error {
	preparedValue, err := types.CounterFromString(value)
	if err != nil {
		c.logger.Warn("Counter converter failed", zap.Error(err))
		return service.ErrInvalidValueFormat
	}

	if err := c.storage.CounterAdd(name, preparedValue); err != nil {
		return fmt.Errorf("add value: %w", err)
	}

	return nil
}

func (c *Counter) Get(name string) (string, error) {
	counterValue, err := c.storage.CounterGet(name)
	if err != nil {
		return "", fmt.Errorf("get clunter value: %w", err)
	}

	return counterValue.String(), nil
}

func (c *Counter) GetAll() (map[string]string, error) {
	counterValues, err := c.storage.GetAllCounterValues()
	if err != nil {
		return nil, fmt.Errorf("get all counter values: %w", err)
	}

	values := make(map[string]string, len(counterValues))

	for name, counterValue := range counterValues {
		values[name] = counterValue.String()
	}

	return values, nil
}

func (c *Counter) GetType() metric.ValueType {
	return metric.Counter
}
