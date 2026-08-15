package config

import (
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

type Agent struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PoolInterval   int    `env:"POOL_INTERVAL"`
	Key            string `env:"KEY"`
}

func ParseAgentConfig() (*Agent, error) {
	agentConfig := &Agent{}

	flag.StringVar(&agentConfig.Address, "a", "localhost:8080", "Server's address and port.")
	flag.IntVar(&agentConfig.PoolInterval, "p", 2, "Metric collection interval.")
	flag.IntVar(&agentConfig.ReportInterval, "r", 10, "Report sending interval.")
	flag.StringVar(&agentConfig.Key, "k", "", "Key for sha256 hashing data.")

	flag.Parse()

	if err := env.Parse(agentConfig); err != nil {
		return nil, fmt.Errorf("parse env variables: %w", err)
	}

	return agentConfig, nil
}
