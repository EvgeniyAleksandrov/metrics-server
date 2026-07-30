//go:generate mockgen -source=processor.go -destination=processor_mock_test.go -package=service_test
package service

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
	models "github.com/EvgeniyAleksandrov/metrics-server/internal/model"
)

type Method interface {
	GetType() string
	Update(update params.Update) error
	Get(name string) (*models.Metrics, error)
	GetAll() (map[string]string, error)
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

func (p *Processor) Update(updateParams params.Update) error {
	method, ok := p.methods[updateParams.MType]
	if !ok {
		return ErrUnsupportedProcessMethod
	}

	if err := method.Update(updateParams); err != nil {
		return fmt.Errorf("update metric: %w", err)
	}

	return nil
}

func (p *Processor) Get(valueParams params.Value) (*models.Metrics, error) {
	method, ok := p.methods[valueParams.MType]
	if !ok {
		return nil, ErrUnsupportedProcessMethod
	}

	modelsMetric, err := method.Get(valueParams.ID)
	if err != nil {
		return nil, fmt.Errorf("get metric: %w", err)
	}

	return modelsMetric, nil
}

func (p *Processor) GetAll() (map[string]map[string]string, error) {
	metrics := make(map[string]map[string]string, len(p.methods))

	for name, method := range p.methods {
		methodMetrics, err := method.GetAll()
		if err != nil {
			return nil, fmt.Errorf("get metrics: %w", err)
		}

		metrics[name] = methodMetrics
	}

	return metrics, nil
}
