package metrics

import "fmt"

type MetricType string

const (
	TypeGauge   = MetricType("gauge")
	TypeCounter = MetricType("counter")
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *int64   `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *float64 `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func NewMetrics(metricType MetricType, metric string, value float64) (*Metrics, error) {
	m := &Metrics{
		ID:    metric,
		MType: string(metricType),
	}

	switch metricType {
	case TypeGauge:
		m.Value = &value
	case TypeCounter:
		v := int64(value)
		m.Delta = &v
	default:
		return nil, fmt.Errorf("Неизвестынй тип метрики %s", metricType)
	}

	return m, nil
}
