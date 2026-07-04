package resource

import (
	"math/rand/v2"
)

const RandomType Type = "random"

type Random struct {
	randValue float64
}

func NewRandom() *Random {
	return &Random{
		randValue: rand.Float64(),
	}
}

func (r *Random) Get() any {
	return r.randValue
}

func (r *Random) Type() Type {
	return RandomType
}

func (r *Random) Update() error {
	r.randValue = rand.Float64()
	return nil
}
