//go:generate mockgen -source=counter.go -destination=counter_mock_test.go -package=method_test
package method

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type CounterStorage interface {
	AddCounter(name string, value int64) error
	GetCounter(name string) (int64, error)
	GetAllCounterValues() (map[string]int64, error)
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

func (c *Counter) Update(updateParams params.Update) error {
	if err := c.storage.AddCounter(updateParams.ID, *updateParams.Delta); err != nil {
		return fmt.Errorf("add counter value: %w", err)
	}

	return nil
}

func (c *Counter) Get(name string) (*models.Metrics, error) {
	counterValue, err := c.storage.GetCounter(name)
	if err != nil {
		return nil, fmt.Errorf("get counter value: %w", err)
	}

	return &models.Metrics{
		ID:    name,
		MType: models.Counter,
		Delta: &counterValue,
	}, nil
}

func (c *Counter) GetAll() (map[string]string, error) {
	counterValues, err := c.storage.GetAllCounterValues()
	if err != nil {
		return nil, fmt.Errorf("get all counter values: %w", err)
	}

	values := make(map[string]string, len(counterValues))

	for name, counterValue := range counterValues {
		values[name] = types.Counter(counterValue).String()
	}

	return values, nil
}

func (c *Counter) GetType() string {
	return models.Counter
}
