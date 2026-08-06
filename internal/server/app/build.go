package app

import (
	"context"
	"fmt"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/config"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/config/db"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/repository"
	"go.uber.org/zap"
)

func BuildStorage(
	ctx context.Context,
	serverConfig *config.Server,
	logger interfaces.Logger,
) (
	interfaces.Storage, error) {
	if serverConfig.DatabaseDSN != "" {
		dbConnection, err := db.BuildDBConnection(ctx, serverConfig.DatabaseDSN)
		if err != nil {
			return nil, fmt.Errorf("build db connection: %w", err)
		}

		return repository.NewDBStorage(dbConnection, logger), nil
	}

	fileStorage := repository.NewStorageSaver(
		ctx,
		repository.NewMemStorage(),
		serverConfig.FilePath, logger,
		serverConfig.StoreInterval,
	)

	if serverConfig.Restore {
		logger.Info("Start data loading")
		if err := fileStorage.Load(); err != nil {
			logger.Warn("Load data from file failed", zap.Error(err))
		} else {
			logger.Info("Data was loaded from file", zap.String("filepath", serverConfig.FilePath))
		}
	}

	return fileStorage, nil
}
