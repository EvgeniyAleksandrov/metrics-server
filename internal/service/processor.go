//go:generate mockgen -source=processor.go -destination=processor_mock_test.go -package=service_test
package service

import (
	"context"
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type Method interface {
	GetType() string
	Update(ctx context.Context, update params.Update) error
	Get(ctx context.Context, name string) (*models.Metrics, error)
	GetAll(ctx context.Context) (map[string]string, error)
}

type Processor struct {
	methods map[string]Method
}

func NewProcessor(methods ...Method) *Processor {
	processor := &Processor{
		methods: make(map[string]Method, len(methods)),
	}

	for _, method := range methods {
		methodType := method.GetType()

		if _, ok := processor.methods[methodType]; ok {
			panic("Try to add already available method")
		}

		processor.methods[methodType] = method
	}

	return processor
}

func (p *Processor) Update(ctx context.Context, updateParams params.Update) error {
	method, ok := p.methods[updateParams.MType]
	if !ok {
		return ErrUnsupportedProcessMethod
	}

	if err := method.Update(ctx, updateParams); err != nil {
		return fmt.Errorf("update metric: %w", err)
	}

	return nil
}

func (p *Processor) Get(ctx context.Context, valueParams params.Value) (*models.Metrics, error) {
	method, ok := p.methods[valueParams.MType]
	if !ok {
		return nil, ErrUnsupportedProcessMethod
	}

	modelsMetric, err := method.Get(ctx, valueParams.ID)
	if err != nil {
		return nil, fmt.Errorf("get metric: %w", err)
	}

	return modelsMetric, nil
}

func (p *Processor) GetAll(ctx context.Context) (map[string]map[string]string, error) {
	metrics := make(map[string]map[string]string, len(p.methods))

	for name, method := range p.methods {
		methodMetrics, err := method.GetAll(ctx)
		if err != nil {
			return nil, fmt.Errorf("get metrics: %w", err)
		}

		metrics[name] = methodMetrics
	}

	return metrics, nil
}
