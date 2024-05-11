package metrics

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"go.uber.org/zap"
)

type StorageData struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]float64 `json:"counters"`
}

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]float64
	mu       sync.RWMutex
	fileMu   sync.RWMutex
	file     *os.File
	logger   zap.SugaredLogger
}

func (s *MemStorage) UpdateGauge(metric string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Gauges[metric] = value
}

func (s *MemStorage) UpdateCounter(metric string, value float64) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.Counters[metric]; !ok {
		s.Counters[metric] = 0
	}

	s.Counters[metric] += value
}

func (s *MemStorage) Updates(values []metrics.Metrics) error {
	for _, metric := range values {
		if metric.MType == string(metrics.TypeGauge) {
			s.UpdateGauge(metric.ID, *metric.Value)
		} else if metric.MType == string(metrics.TypeCounter) {
			s.UpdateCounter(metric.ID, float64(*metric.Delta))
		} else {
			return fmt.Errorf("неизвестный тип метрики %v", metric.MType)
		}
	}

	return nil
}

func (s *MemStorage) GetGauges() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]float64)
	for k, v := range s.Gauges {
		result[k] = v
	}

	return result
}

func (s *MemStorage) GetCounters() map[string]float64 {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make(map[string]float64)
	for k, v := range s.Counters {
		result[k] = v
	}

	return result
}

func (s *MemStorage) GetValue(metricType metrics.MetricType, metric string) (float64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if metricType == metrics.TypeGauge {
		v, ok := s.Gauges[metric]

		return v, ok
	}

	if metricType == metrics.TypeCounter {
		v, ok := s.Counters[metric]

		return v, ok
	}

	return 0, false
}

func (s *MemStorage) Store() error {
	s.logger.Infoln("Start store to file - ", s.file.Name())

	if s.file == nil {
		s.logger.Infoln("No file for store")
		return nil
	}

	s.fileMu.Lock()
	defer s.fileMu.Unlock()

	data := &StorageData{
		Gauges:   s.Gauges,
		Counters: s.Counters,
	}

	dataJSON, err := json.Marshal(data)

	if err != nil {
		return err
	}

	dataJSON = append(dataJSON, '\n')

	s.logger.Infoln("Data for store", string(dataJSON))

	err = clearFileAndWriteIntoHim(s.file, dataJSON)

	if err != nil {
		return err
	}

	s.logger.Infoln("Store to file success")

	return nil
}

func (s *MemStorage) Restore() error {
	s.logger.Infoln("Start restore from file ", s.file.Name())

	if s.file == nil {
		s.logger.Infoln("No file for restore")
		return nil
	}

	s.fileMu.Lock()
	defer s.fileMu.Unlock()

	reader := bufio.NewReader(s.file)
	b, err := reader.ReadBytes('\n')

	// если файл просто пустой, то это нормально, просто нет данных для восстановления
	if err == io.EOF {
		s.logger.Infoln("Empty file for restore")
		return nil
	} else if err != nil {
		return err
	}

	s.logger.Infoln("Loaded data for restore ", string(b))

	data := StorageData{}
	err = json.Unmarshal(b, &data)

	if err != nil {
		return err
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	s.Gauges = data.Gauges
	s.Counters = data.Counters

	s.logger.Infoln("End restore from file")

	return nil
}

func (s *MemStorage) Close() error {
	if s.file == nil {
		return nil
	}

	s.logger.Infoln("Close file")

	return s.file.Close()
}

func (s *MemStorage) StoreAndClose() error {
	s.logger.Infoln("StoreAndClose")
	err := s.Store()
	if err != nil {
		return err
	}
	return s.Close()
}

func NewMemStorage(filePath string, logger zap.SugaredLogger) (*MemStorage, error) {
	file, err := openFile(filePath)

	if err != nil {
		return nil, err
	}

	return &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]float64),
		file:     file,
		logger:   logger,
	}, nil
}

func openFile(filePath string) (*os.File, error) {
	if len(filePath) == 0 {
		return nil, nil
	}

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0644)

	if err != nil {
		return nil, err
	}

	return file, nil
}

func clearFileAndWriteIntoHim(f *os.File, b []byte) error {
	err := f.Truncate(0)

	if err != nil {
		return err
	}

	_, err = f.Seek(0, 0)

	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}
