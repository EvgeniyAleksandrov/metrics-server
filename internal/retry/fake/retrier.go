package fake

import "context"

type NopeRetrier struct{}

func NewNopeRetrier() *NopeRetrier {
	return &NopeRetrier{}
}

func (r *NopeRetrier) Do(_ context.Context, fn func() error) error {
	return fn()
}
