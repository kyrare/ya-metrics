package metrics

import "github.com/kyrare/ya-metrics/internal/domain/metrics"

type Storage interface {
	UpdateGauge(metric string, value float64)
	UpdateCounter(metric string, value float64)
	GetGauges() map[string]float64
	GetCounters() map[string]float64
	GetValue(metricType metrics.MetricType, metric string) (float64, bool)
	Store() error
	Restore() error
	Close() error
	StoreAndClose() error
}
