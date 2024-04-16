package metrics

import (
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"sync"
)

type Storage interface {
	UpdateGauge(metric string, value float64)
	UpdateCounter(metric string, value float64)
	GetGauges() map[string]float64
	GetCounters() map[string]float64
	Get(metricType metrics.MetricType, metric string) (float64, bool)
}

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string]float64
	mu       sync.RWMutex
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

func NewMemStorage() *MemStorage {
	return &MemStorage{
		Gauges:   make(map[string]float64),
		Counters: make(map[string]float64),
	}
}
