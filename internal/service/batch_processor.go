package service

import (
	"context"
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/params"
)

type BatchMethod interface {
	GetType() string
	UpdateBatch(ctx context.Context, updates []params.Update) error
}

type BatchProcessor struct {
	batchMethods map[string]BatchMethod
}

func NewBatchProcessor(batchMethods ...BatchMethod) *BatchProcessor {
	processor := &BatchProcessor{
		batchMethods: make(map[string]BatchMethod, len(batchMethods)),
	}

	for _, batchMethod := range batchMethods {
		methodType := batchMethod.GetType()

		if _, ok := processor.batchMethods[methodType]; ok {
			panic("Try to add already available method")
		}

		processor.batchMethods[methodType] = batchMethod
	}

	return processor
}

func (p *BatchProcessor) UpdateBatch(ctx context.Context, updates []params.Update) error {
	updateData := make(map[string][]params.Update)

	for _, update := range updates {
		mType := update.MType

		_, ok := p.batchMethods[mType]
		if !ok {
			return ErrUnsupportedProcessMethod
		}

		if _, ok := updateData[mType]; !ok {
			updateData[mType] = make([]params.Update, 0)
		}

		updateData[mType] = append(updateData[mType], update)
	}

	for mType, mUpdates := range updateData {
		if err := p.batchMethods[mType].UpdateBatch(ctx, mUpdates); err != nil {
			return fmt.Errorf("update batch method '%s': %w", mType, err)
		}
	}

	return nil
}
