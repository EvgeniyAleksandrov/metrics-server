package main

import (
	"flag"
	"net/http"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/config"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/middlware"
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
	memStorage := repository.NewMemStorage()

	router := chi.NewRouter()

	processor := service.NewProcessor(method.NewGauge(memStorage, logger), method.NewCounter(memStorage, logger))

	router.Get("/", middlware.WithLogging(handler.NewRoot(processor).ServeHTTP, logger))

	router.Post(
		"/value/",
		middlware.WithLogging(
			handler.NewValue(processor, parservalue.NewJSON(MaxRequestSize), value.NewJSONWriter()).ServeHTTP,
			logger,
		),
	)
	router.Post(
		"/value",
		middlware.WithLogging(
			handler.NewValue(processor, parservalue.NewJSON(MaxRequestSize), value.NewJSONWriter()).ServeHTTP,
			logger,
		),
	)

	router.Post(
		"/update/",
		middlware.WithLogging(
			handler.NewUpdate(processor, parserupdate.NewJSON(MaxRequestSize), update.NewJSONWriter()).ServeHTTP,
			logger,
		),
	)

	router.Post(
		"/update",
		middlware.WithLogging(
			handler.NewUpdate(processor, parserupdate.NewJSON(MaxRequestSize), update.NewJSONWriter()).ServeHTTP,
			logger,
		),
	)

	router.Post(
		"/update/{method}/{name}/{value}",
		middlware.WithLogging(
			handler.NewUpdate(processor, parserupdate.NewPath(logger), update.NewValueWriter()).ServeHTTP,
			logger,
		),
	)

	router.Get(
		"/value/{method}/{name}",
		middlware.WithLogging(
			handler.NewValue(processor, parservalue.NewPath(), value.NewValueWriter()).ServeHTTP,
			logger,
		),
	)

	return http.ListenAndServe(addr, router)
}
