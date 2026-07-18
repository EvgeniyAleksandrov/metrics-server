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
	value := c.pull
	c.pull = 0

	return value
}

func (c *PullCounter) Update() error {
	c.pull++
	return nil
}
