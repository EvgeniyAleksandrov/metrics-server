package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/agent"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/config"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/getter"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
	"github.com/caarlos0/env"
	"go.uber.org/zap"
)

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(fmt.Sprintf("Can't create logger: %s", err.Error()))
	}

	defer logger.Sync()

	agentConfig := config.Agent{}

	flag.StringVar(&agentConfig.Address, "a", "localhost:8080", "Server's address and port.")
	flag.IntVar(&agentConfig.PoolInterval, "p", 2, "Metric collection interval.")
	flag.IntVar(&agentConfig.ReportInterval, "r", 10, "Report sending interval.")

	flag.Parse()

	if err := env.Parse(&agentConfig); err != nil {
		logger.Fatal("Parse config error", zap.Error(err))
	}

	a := agent.NewAgent(
		resource.NewManager(
			resource.NewMemory(),
			resource.NewPullCounter(),
			resource.NewRandom(),
		),
		[]agent.MetricGetter{
			getter.NewRandom(),
			getter.NewPullCounter(),
			getter.NewMemory(),
		},
		metric.NewPublisher(&http.Client{}),
		agentConfig.Address,
		time.Duration(agentConfig.PoolInterval)*time.Second,
		time.Duration(agentConfig.ReportInterval)*time.Second,
		logger,
	)

	stopContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	logger.Info(
		"Agent started",
		zap.String("addr", agentConfig.Address),
		zap.Int("pollInterval", agentConfig.PoolInterval),
		zap.Int("reportInterval", agentConfig.ReportInterval),
	)

	go func() {
		a.Run(stopContext)
		wg.Done()
	}()

	wg.Wait()

	logger.Info("Agent finished")
}
