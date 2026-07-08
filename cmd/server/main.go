package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/method"
	"github.com/go-chi/chi/v5"
)

func main() {
	var addr string

	flag.StringVar(
		&addr,
		"a",
		"localhost:8080",
		"The address and port on which the server listens for connections.",
	)

	flag.Parse()

	log.Printf("Start server on: %s", addr)

	if err := run(addr); err != nil {
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
