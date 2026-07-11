package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/agent"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/getter"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
)

func main() {
	var (
		serverAddr     string
		pullInterval   int
		reportInterval int
	)

	flag.StringVar(&serverAddr, "a", "localhost:8080", "Server's address and port.")
	flag.IntVar(&pullInterval, "p", 2, "Metric collection interval.")
	flag.IntVar(&reportInterval, "r", 10, "Report sending interval.")

	flag.Parse()

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
		serverAddr,
		time.Duration(pullInterval)*time.Second,
		time.Duration(reportInterval)*time.Second,
	)

	stopContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	var wg sync.WaitGroup
	wg.Add(1)

	log.Println("Agent started.")

	go func() {
		a.Run(stopContext)
		wg.Done()
	}()

	wg.Wait()

	log.Println("Agent stoped.")

}
