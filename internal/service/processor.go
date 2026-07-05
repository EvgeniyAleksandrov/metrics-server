//go:generate mockgen -source=processor.go -destination=processor_mock_test.go -package=service_test
package service

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
)

type Method interface {
	GetType() metric.ValueType
	Update(name, value string) error
	Get(name string) (string, error)
	GetAll() (map[string]string, error)
}

type Processor struct {
	methods map[metric.ValueType]Method
}

func NewProcessor(methods ...Method) *Processor {
	processor := &Processor{
		methods: make(map[metric.ValueType]Method, len(methods)),
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

func (p *Processor) Update(methodType, name, value string) error {
	method, ok := p.methods[metric.ValueType(methodType)]
	if !ok {
		return ErrUnsupportedProcessMethod
	}

	if err := method.Update(name, value); err != nil {
		return fmt.Errorf("update metric: %w", err)
	}

	return nil
}

func (p *Processor) Get(methodType string, name string) (string, error) {
	method, ok := p.methods[metric.ValueType(methodType)]
	if !ok {
		return "", ErrUnsupportedProcessMethod
	}

	value, err := method.Get(name)
	if err != nil {
		return "", fmt.Errorf("get metric: %w", err)
	}

	return value, nil
}

func (p *Processor) GetAll() (map[string]map[string]string, error) {
	metrics := make(map[string]map[string]string, len(p.methods))

	for name, method := range p.methods {
		methodMetrics, err := method.GetAll()
		if err != nil {
			return nil, fmt.Errorf("get metrics: %w", err)
		}

		metrics[string(name)] = methodMetrics
	}

	return metrics, nil
}
