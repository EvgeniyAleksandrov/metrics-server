package main

import (
	"flag"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/config"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/middleware"
	parserupdate "github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/parser/update"
	parservalue "github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/parser/value"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/response/update"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/response/value"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/method"
	"github.com/caarlos0/env"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

const MaxRequestSize = 8 * 1024 * 1024

func main() {
	logger, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	serverConfig := config.Server{}

	flag.StringVar(
		&serverConfig.Address,
		"a",
		"localhost:8080",
		"The address and port on which the server listens for connections.",
	)

	flag.Parse()

	if err := env.Parse(&serverConfig); err != nil {
		logger.Fatal("Parse config error.", zap.Error(err))
	}

	logger.Info("Start server. ", zap.String("address", serverConfig.Address))

	if err := run(serverConfig.Address, logger); err != nil {
		logger.Sugar().Fatalf("server run: %s", err.Error())
	}
}

func run(addr string, logger *zap.Logger) error {
	// Build Services
	memStorage := repository.NewMemStorage()
	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	// Build Handlers
	jsonValueHandler := handler.NewValue(processor, parservalue.NewJSON(MaxRequestSize), value.NewJSONWriter())
	pathValueHandler := handler.NewValue(processor, parservalue.NewPath(), value.NewValueWriter())
	jsonUpdateHandler := handler.NewUpdate(processor, parserupdate.NewJSON(MaxRequestSize), update.NewJSONWriter())
	pathUpdateHandler := handler.NewUpdate(processor, parserupdate.NewPath(logger), update.NewValueWriter())

	// Build router
	router := chi.NewRouter()

	// middlewares
	router.Use(middleware.NewCompression(logger).Do, middleware.NewLogging(logger).Do)

	// handlers
	router.Get("/", handler.NewRoot(processor).ServeHTTP)

	router.Post("/value/", jsonValueHandler.ServeHTTP)
	router.Post("/value", jsonValueHandler.ServeHTTP)
	router.Post("/update/", jsonUpdateHandler.ServeHTTP)
	router.Post("/update", jsonUpdateHandler.ServeHTTP)

	router.Post("/update/{method}/{name}/{value}", pathUpdateHandler.ServeHTTP)
	router.Get("/value/{method}/{name}", pathValueHandler.ServeHTTP)

	return http.ListenAndServe(addr, router)
}
