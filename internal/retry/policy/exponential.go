package policy

import "time"

//type Policy interface {
//	ShouldRetry(err error) bool
//	Delay(attempt int) *time.Duration
//}

type ExponentialBackoff struct {
	maxAttempt int
	step       int
}

func NewExponentialBackOff(maxAttempt, step int) *ExponentialBackoff {
	return &ExponentialBackoff{
		maxAttempt: maxAttempt,
		step:       step,
	}
}

func (b *ExponentialBackoff) ShouldRetry(err error) bool {
	if err == nil {
		return false
	}

	return true
}

func (b *ExponentialBackoff) Delay(attempt int) *time.Duration {
	if attempt > b.maxAttempt {
		return nil
	}

	delay := time.Duration(b.step)*time.Duration(attempt)*time.Second - time.Second

	return &delay
}
