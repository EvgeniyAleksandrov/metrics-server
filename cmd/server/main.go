package main

import (
	"log"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/update"
	method2 "github.com/EvgeniyAleksandrov/metrics-server/internal/service/update/method"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server run: %s", err.Error())
	}
}

func run() error {
	memStorage := repository.NewMemStorage()

	mux := http.NewServeMux()

	mux.Handle(
		"/update/{method}/{name}/{value}",
		handler.NewUpdate(
			update.NewProcessor(
				method2.NewGauge(memStorage),
				method2.NewCounter(memStorage),
			),
		),
	)

	return http.ListenAndServe("localhost:8080", mux)
}
