package update

import (
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/metric"
)

type UpdateMethod interface {
	GetType() string
	Update(name, value string) error
}

type Processor struct {
	methods map[string]UpdateMethod
}

func NewProcessor(methods ...UpdateMethod) *Processor {
	processor := &Processor{
		methods: make(map[string]UpdateMethod, len(methods)),
	}

	for _, method := range methods {
		if _, ok := processor.methods[method.GetType()]; ok {
			panic("Try to add available method")
		}

		processor.methods[method.GetType()] = method
	}

	return processor
}

func (p *Processor) Update(methodType, name, value string) error {
	method, ok := p.methods[methodType]
	if !ok {
		return metric.ErrUnsupportedProcessMethod
	}

	if err := method.Update(name, value); err != nil {
		return fmt.Errorf("update metric: %w", err)
	}

	return nil
}
