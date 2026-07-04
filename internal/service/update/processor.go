//go:generate mockgen -source=processor.go -destination=processor_mock_test.go -package=update_test
package update

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
)

type Method interface {
	GetType() metric.ValueType
	Update(name, value string) error
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
