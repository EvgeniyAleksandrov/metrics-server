package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/config"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/method"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
)

func main() {
	serverConfig := config.Server{}

	flag.StringVar(
		&serverConfig.Address,
		"a",
		"localhost:8080",
		"The address and port on which the server listens for connections.",
	)

	flag.Parse()

	if err := env.Parse(&serverConfig); err != nil {
		log.Fatalf("Parse config error: %s", err.Error())
	}

	log.Printf("Start server on: %s", serverConfig.Address)

	if err := run(serverConfig.Address); err != nil {
		log.Fatalf("server run: %s", err.Error())
	}
}

func run(addr string) error {
	memStorage := repository.NewMemStorage()

	router := chi.NewRouter()

	processor := service.NewProcessor(method.NewGauge(memStorage), method.NewCounter(memStorage))

	router.Get("/", handler.NewRoot(processor).ServeHTTP)
	router.Get("/value/{method}/{name}", handler.NewGetValue(processor).ServeHTTP)
	router.Post("/update/{method}/{name}/{value}", handler.NewUpdate(processor).ServeHTTP)

	return http.ListenAndServe(addr, router)
}
