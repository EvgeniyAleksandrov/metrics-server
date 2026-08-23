package config

import (
	"errors"
	"flag"
	"fmt"

	"github.com/caarlos0/env"
)

var (
	ErrRateLimiterTooLow           = errors.New("too low num of rate limiter")
	ErrPoolIntervalTooLow          = errors.New("too low pool interval")
	ErrReportIntervalTooLow        = errors.New("too low Report interval")
	ErrReportIntervalLowerThanPool = errors.New("report interval lower then pool interval")
)

type Agent struct {
	Address        string `env:"ADDRESS"`
	ReportInterval int    `env:"REPORT_INTERVAL"`
	PoolInterval   int    `env:"POOL_INTERVAL"`
	Key            string `env:"KEY"`
	RateLimit      int    `env:"RATE_LIMIT"`
}

func ParseAgentConfig() (*Agent, error) {
	agentConfig := &Agent{}

	flag.StringVar(&agentConfig.Address, "a", "localhost:8080", "Server's address and port.")
	flag.IntVar(&agentConfig.PoolInterval, "p", 2, "Metric collection interval.")
	flag.IntVar(&agentConfig.ReportInterval, "r", 10, "Report sending interval.")
	flag.StringVar(&agentConfig.Key, "k", "", "Key for sha256 hashing data.")
	flag.IntVar(&agentConfig.RateLimit, "l", 1, "Number of send metric workers.")

	flag.Parse()

	if err := env.Parse(agentConfig); err != nil {
		return nil, fmt.Errorf("parse env variables: %w", err)
	}

	validationErrors := make([]error, 0)

	if agentConfig.RateLimit < 1 {
		validationErrors = append(validationErrors, ErrRateLimiterTooLow)
	}

	if agentConfig.PoolInterval < 1 {
		validationErrors = append(validationErrors, ErrPoolIntervalTooLow)
	}

	if agentConfig.ReportInterval < 1 {
		validationErrors = append(validationErrors, ErrReportIntervalTooLow)
	}

	if agentConfig.PoolInterval > agentConfig.ReportInterval {
		validationErrors = append(validationErrors, ErrReportIntervalLowerThanPool)
	}

	return agentConfig, errors.Join(validationErrors...)
}
