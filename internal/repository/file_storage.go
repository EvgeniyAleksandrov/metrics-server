package repository

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
)

type Logger interface {
	Warn(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
}

type StorageSaver struct {
	filePath     string
	storage      *MemStorage
	logger       Logger
	saveInterval int

	wg sync.WaitGroup
}

func NewStorageSaver(ctx context.Context, storage *MemStorage, filePath string, logger Logger, saveInterval int) *StorageSaver {
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
			if err := s.Save(); err != nil {
				s.logger.Warn("Save storage failed", zap.Error(err))
			}
			s.logger.Info("Stop storage write process")
			return
		}
	}
}

func (s *StorageSaver) SetGauge(name string, value float64) error {
	err := s.storage.SetGauge(name, value)
	if err != nil {
		return err
	}

	// Синхронное сохранение, если интервал равен 0
	if s.saveInterval == 0 {
		return s.Save()
	}

	return nil
}

func (s *StorageSaver) AddCounter(name string, value int64) error {
	err := s.storage.AddCounter(name, value)
	if err != nil {
		return err
	}

	if s.saveInterval == 0 {
		return s.Save()
	}

	return nil
}

func (s *StorageSaver) GetGauge(name string) (float64, error) {
	return s.storage.GetGauge(name)
}

func (s *StorageSaver) GetCounter(name string) (int64, error) {
	return s.storage.GetCounter(name)
}

func (s *StorageSaver) GetAllGaugeValues() (map[string]float64, error) {
	return s.storage.GetAllGaugeValues()
}

func (s *StorageSaver) GetAllCounterValues() (map[string]int64, error) {
	return s.storage.GetAllCounterValues()
}
