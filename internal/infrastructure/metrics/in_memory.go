package metrics

import (
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
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
}

func (storage *MemStorage) UpdateGauge(metric string, value float64) {
	storage.Gauges[metric] = value
}

func (storage *MemStorage) UpdateCounter(metric string, value float64) {
	if _, ok := storage.Counters[metric]; !ok {
		storage.Counters[metric] = 0
	}

	storage.Counters[metric] += value
}

func (storage *MemStorage) GetGauges() map[string]float64 {
	return storage.Gauges
}

func (storage *MemStorage) GetCounters() map[string]float64 {
	return storage.Counters
}

func (storage MemStorage) Get(metricType metrics.MetricType, metric string) (float64, bool) {
	if metricType == metrics.TypeGauge {
		v, ok := storage.Gauges[metric]

		return v, ok
	}

	if metricType == metrics.TypeCounter {
		v, ok := storage.Counters[metric]

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
