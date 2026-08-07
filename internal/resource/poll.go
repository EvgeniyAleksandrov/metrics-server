package resource

const PullCounterType Type = "pollCount-counter"

type PollCounter struct {
	pollCount int64
}

func NewPollCounter() *PollCounter {
	return &PollCounter{}
}

func (c *PollCounter) Type() Type {
	return PullCounterType
}

func (c *PollCounter) Get() any {
	value := c.pollCount
	c.pollCount = 0

	return value
}

func (c *PollCounter) Update() error {
	c.pollCount++
	return nil
}
