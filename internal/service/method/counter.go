package method

import (
	"context"
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/types"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type Counter struct {
	storage interfaces.Storage
	logger  interfaces.Logger
}

func NewCounter(storage interfaces.Storage, logger interfaces.Logger) *Counter {
	return &Counter{
		storage: storage,
		logger:  logger,
	}
}

func (c *Counter) Update(ctx context.Context, updateParams params.Update) error {
	if err := c.storage.AddCounter(ctx, updateParams.ID, *updateParams.Delta); err != nil {
		return fmt.Errorf("add counter value: %w", err)
	}

	return nil
}

func (c *Counter) Get(ctx context.Context, name string) (*models.Metrics, error) {
	counterValue, err := c.storage.GetCounter(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("get counter value: %w", err)
	}

	return &models.Metrics{
		ID:    name,
		MType: models.Counter,
		Delta: &counterValue,
	}, nil
}

func (c *Counter) GetAll(ctx context.Context) (map[string]string, error) {
	counterValues, err := c.storage.GetAllCounterValues(ctx)
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
