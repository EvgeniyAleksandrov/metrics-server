package resource

import (
	"errors"
	"fmt"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
)

const SystemType Type = "system"

var ErrIncorrectPercentOfCPU = errors.New("incorrect cpu percent")

type SystemStats struct {
	MemTotal   float64
	FreeMemory float64
	CPUUsage   []float64
}

type System struct {
	stats SystemStats
}

func NewSystem() *System {
	return &System{
		stats: SystemStats{},
	}
}

func (s *System) Get() any { return s.stats }

func (s *System) Type() Type { return SystemType }

func (s *System) Update() error {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return fmt.Errorf("get virtual memory: %w", err)
	}

	cpuUsage, err := cpu.Percent(1*time.Second, true)
	if err != nil {
		return fmt.Errorf("get cpu percent usage: %w", err)
	}

	s.stats.MemTotal = float64(memory.Total)
	s.stats.FreeMemory = float64(memory.Free)
	s.stats.CPUUsage = cpuUsage

	return nil
}
