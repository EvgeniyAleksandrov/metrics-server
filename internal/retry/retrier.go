package retry

import (
	"context"
	"time"
)

type ErrorChecker interface {
	ShouldRetry(err error) bool
}

type Policy interface {
	Delay(attempt int) *time.Duration
}

type Retry struct {
	errorChecker ErrorChecker
	policy       Policy
}

func NewRetry(errorChecker ErrorChecker, policy Policy) *Retry {
	return &Retry{
		errorChecker: errorChecker,
		policy:       policy,
	}
}

func (r Retry) Do(ctx context.Context, fn func() error) error {
	for attempt := 1; ; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}

		if !r.errorChecker.ShouldRetry(err) {
			return err
		}

		delay := r.policy.Delay(attempt)
		if delay == nil {
			return err
		}

		select {
		case <-time.After(*delay):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
