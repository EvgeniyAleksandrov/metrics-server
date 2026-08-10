package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/config"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/middleware"
	parserbatchupdate "github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/parser/batch/update"
	parserupdate "github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/parser/update"
	parservalue "github.com/EvgeniyAleksandrov/metrics-server/internal/handler/request/parser/value"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/response/update"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/handler/response/value"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/server/app"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/service/method"
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

	serverConfig, err := config.ParseServerConfig()
	if err != nil {
		logger.Fatal("Parse config error.", zap.Error(err))
	}

	stopContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT, syscall.SIGKILL)
	defer stop()

	if err := run(stopContext, serverConfig, logger); err != nil {
		logger.Sugar().Fatalf("server run: %s", err.Error())
	}
}

func run(ctx context.Context, serverConfig *config.Server, logger *zap.Logger) error {
	// Build Services
	storage, err := app.BuildStorage(ctx, serverConfig, logger)
	if err != nil {
		return fmt.Errorf("build storage: %w", err)
	}

	defer storage.Close()

	counterMethod := method.NewCounter(storage, logger)
	gaugeMethod := method.NewGauge(storage, logger)

	processor := service.NewProcessor(gaugeMethod, counterMethod)
	batchProcessor := service.NewBatchProcessor(gaugeMethod, counterMethod)

	// Build Handlers
	jsonValueHandler := handler.NewValue(processor, parservalue.NewJSON(MaxRequestSize), value.NewJSONWriter())
	pathValueHandler := handler.NewValue(processor, parservalue.NewPath(), value.NewValueWriter())
	jsonUpdateHandler := handler.NewUpdate(processor, parserupdate.NewJSON(MaxRequestSize), update.NewJSONWriter())
	pathUpdateHandler := handler.NewUpdate(processor, parserupdate.NewPath(logger), update.NewValueWriter())

	jsonUpdateBatchHandler := handler.NewUpdateBatch(
		batchProcessor,
		parserbatchupdate.NewJSON(MaxRequestSize),
		update.NewJSONWriter(),
		logger,
	)

	// Build router
	router := chi.NewRouter()

	// middlewares
	router.Use(middleware.NewCompression(logger).Do, middleware.NewLogging(logger).Do)

	// handlers
	router.Get("/", handler.NewRoot(processor).ServeHTTP)
	router.Get("/ping", handler.NewPing(storage, logger).ServeHTTP)
	router.Get("/ping/", handler.NewPing(storage, logger).ServeHTTP)

	router.Post("/value/", jsonValueHandler.ServeHTTP)
	router.Post("/value", jsonValueHandler.ServeHTTP)
	router.Post("/update/", jsonUpdateHandler.ServeHTTP)
	router.Post("/update", jsonUpdateHandler.ServeHTTP)
	router.Post("/updates/", jsonUpdateBatchHandler.ServeHTTP)
	router.Post("/updates", jsonUpdateBatchHandler.ServeHTTP)

	router.Post("/update/{method}/{name}/{value}", pathUpdateHandler.ServeHTTP)
	router.Get("/value/{method}/{name}", pathValueHandler.ServeHTTP)

	srv := &http.Server{
		Addr:    serverConfig.Address,
		Handler: router,
	}

	errChan := make(chan error, 1)

	go func() {
		logger.Info("Server is started", zap.String("address", serverConfig.Address))
		errChan <- srv.ListenAndServe()
	}()

	select {
	case err := <-errChan:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	case <-ctx.Done():
		logger.Info("Shutting down server gracefully...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// Корректно останавливаем сервер
		return srv.Shutdown(shutdownCtx)
	}
}
