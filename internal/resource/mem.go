package resource

import (
	"runtime"
)

const MemoryType Type = "memory"

type Memory struct {
	memStat runtime.MemStats
}

func NewMemory() *Memory {
	return &Memory{
		memStat: runtime.MemStats{},
	}
}

func (m *Memory) Type() Type {
	return MemoryType
}

func (m *Memory) Get() any {
	return m.memStat
}

func (m *Memory) Update() error {
	runtime.ReadMemStats(&m.memStat)
	return nil
}
