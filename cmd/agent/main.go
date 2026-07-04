package main

import (
	"log"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/agent"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/getter"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/resource"
)

func main() {

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
		"localhost:8080",
	)

	if err := a.Run(); err != nil {
		log.Fatalf("Error: %s", err.Error())
	}

}
