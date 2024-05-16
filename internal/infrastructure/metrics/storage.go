package metrics

import (
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
)

type storageGauge interface {
	UpdateGauge(metric string, value float64)
	GetGauges() map[string]float64
}

type storageCounter interface {
	UpdateCounter(metric string, value float64)
	GetCounters() map[string]float64
}

type storageStored interface {
	Store() error
	Restore() error
	StoreAndClose() error
}

type Storage interface {
	storageGauge
	storageCounter
	storageStored
	GetValue(metricType metrics.MetricType, metric string) (float64, bool)
	Updates(values []metrics.Metrics) error
	Close() error
}
