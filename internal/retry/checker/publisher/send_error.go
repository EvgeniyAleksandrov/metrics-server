package publisher

import (
	"errors"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
)

type SendError struct{}

func NewSendError() *SendError {
	return &SendError{}
}

func (c *SendError) ShouldRetry(err error) bool {
	return errors.Is(err, metric.ErrNotPublished)
}
