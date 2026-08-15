package repository

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/EvgeniyAleksandrov/metrics-server/internal/interfaces"
	"github.com/EvgeniyAleksandrov/metrics-server/internal/metric/repository/values"
	"go.uber.org/zap"
)

type StorageSaver struct {
	filePath     string
	storage      interfaces.MarshaledStorage
	logger       interfaces.Logger
	saveInterval int

	wg sync.WaitGroup
}

func NewStorageSaver(
	ctx context.Context,
	storage interfaces.MarshaledStorage,
	filePath string,
	logger interfaces.Logger,
	saveInterval int,
) *StorageSaver {
	fileStorage := &StorageSaver{
		filePath:     filePath,
		storage:      storage,
		logger:       logger,
		saveInterval: saveInterval,
	}

	if saveInterval > 0 {
		fileStorage.wg.Add(1)
		go fileStorage.RunStorageSaveProcess(ctx)
	}

	return fileStorage
}

func (s *StorageSaver) Save() error {
	storageData, err := s.storage.Marshal()
	if err != nil {
		return fmt.Errorf("marshal storage data: %w", err)
	}

	if err := os.WriteFile(s.filePath, storageData, 0644); err != nil {
		return fmt.Errorf("write data to file: %w", err)
	}

	return nil
}

func (s *StorageSaver) Load() error {
	fileData, err := os.ReadFile(s.filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	s.logger.Info("File loaded")
	if err := s.storage.Unmarshal(fileData); err != nil {
		return fmt.Errorf("unmarshal file data: %w", err)
	}

	s.logger.Info("File was unmarshaled")
	return nil
}

func (s *StorageSaver) RunStorageSaveProcess(ctx context.Context) {
	if s.saveInterval <= 0 {
		return
	}

	defer s.wg.Done()

	saveTicker := time.NewTicker(time.Duration(s.saveInterval) * time.Second)
	defer saveTicker.Stop()

	for {
		select {
		case <-saveTicker.C:
			if err := s.Save(); err != nil {
				s.logger.Warn("Save storage failed", zap.Error(err))
				continue
			}

			s.logger.Info("Storage was saved", zap.String("filePath", s.filePath))
		case <-ctx.Done():
			s.logger.Info("Stop storage write process")
			return
		}
	}
}

func (s *StorageSaver) SetGauge(ctx context.Context, gauge values.Gauge) error {
	err := s.storage.SetGauge(ctx, gauge)
	if err != nil {
		return err
	}

	// Синхронное сохранение, если интервал равен 0
	if s.saveInterval == 0 {
		return s.Save()
	}

	return nil
}

func (s *StorageSaver) SetGauges(ctx context.Context, gauges []values.Gauge) error {
	err := s.storage.SetGauges(ctx, gauges)
	if err != nil {
		return err
	}

	// Синхронное сохранение, если интервал равен 0
	if s.saveInterval == 0 {
		return s.Save()
	}

	return nil
}

func (s *StorageSaver) AddCounter(ctx context.Context, counter values.Counter) error {
	err := s.storage.AddCounter(ctx, counter)
	if err != nil {
		return err
	}

	if s.saveInterval == 0 {
		return s.Save()
	}

	return nil
}

func (s *StorageSaver) AddCounters(ctx context.Context, counters []values.Counter) error {
	err := s.storage.AddCounters(ctx, counters)
	if err != nil {
		return err
	}

	if s.saveInterval == 0 {
		return s.Save()
	}

	return nil
}

func (s *StorageSaver) GetGauge(ctx context.Context, name string) (float64, error) {
	return s.storage.GetGauge(ctx, name)
}

func (s *StorageSaver) GetCounter(ctx context.Context, name string) (int64, error) {
	return s.storage.GetCounter(ctx, name)
}

func (s *StorageSaver) GetAllGaugeValues(ctx context.Context) (map[string]float64, error) {
	return s.storage.GetAllGaugeValues(ctx)
}

func (s *StorageSaver) GetAllCounterValues(ctx context.Context) (map[string]int64, error) {
	return s.storage.GetAllCounterValues(ctx)
}

func (s *StorageSaver) Ping(_ context.Context) error {
	return nil
}

func (s *StorageSaver) Close() {
	s.logger.Info("Close file storage with save data")
	if err := s.Save(); err != nil {
		s.logger.Warn("Close storage with save error", zap.Error(err))
	}
}
