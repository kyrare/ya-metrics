package metrics

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
)

type Storage interface {
	UpdateGauge(metric string, value float64)
	UpdateCounter(metric string, value float64)
	GetGauges() map[string]float64
	GetCounters() map[string]float64
	Get(metricType metrics.MetricType, metric string) (float64, bool)
	Store() error
	Restore() error
}

type StorageData struct {
	Gauges   map[string]float64 `json:"gauges"`
	Counters map[string]float64 `json:"counters"`
}

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]float64
	mu       sync.RWMutex
	fileMu   sync.RWMutex
	filePath string
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

func (s *MemStorage) Get(metricType metrics.MetricType, metric string) (float64, bool) {
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
	if len(s.filePath) == 0 {
		return nil
	}

	data := &StorageData{
		Gauges:   s.Gauges,
		Counters: s.Counters,
	}

	dataJSON, err := json.Marshal(data)

	if err != nil {
		return err
	}

	dataJSON = append(dataJSON, '\n')

	s.fileMu.Lock()
	defer s.fileMu.Unlock()

	file, err := os.OpenFile(s.filePath, os.O_WRONLY|os.O_CREATE, 0644)
	defer func() {
		err := file.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	if err != nil {
		return err
	}

	_, err = file.Write(dataJSON)

	return err
}

func (s *MemStorage) Restore() error {
	if len(s.filePath) == 0 {
		return nil
	}

	// проверяем, что файл существует
	if _, err := os.Stat(s.filePath); errors.Is(err, os.ErrNotExist) {
		return nil
	}

	s.fileMu.RLock()
	defer s.fileMu.RUnlock()

	file, err := os.Open(s.filePath)

	if err != nil {
		return err
	}

	reader := bufio.NewReader(file)
	b, err := reader.ReadBytes('\n')

	if err != nil {
		return err
	}

	data := StorageData{}
	err = json.Unmarshal(b, &data)

	if err != nil {
		return err
	}

	s.Gauges = data.Gauges
	s.Counters = data.Counters

	return nil
}

func NewMemStorage(filePath string) *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]float64),
		filePath: filePath,
	}
}
