package resource

const PullCounterType Type = "pull-counter"

type PullCounter struct {
	pull int64
}

func NewPullCounter() *PullCounter {
	return &PullCounter{}
}

func (c *PullCounter) Type() Type {
	return PullCounterType
}

func (c *PullCounter) Get() any {
	return c.pull
}

func (c *PullCounter) Update() error {
	c.pull++
	return nil
}
