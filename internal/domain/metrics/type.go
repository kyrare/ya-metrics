package metrics

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

func NewMetrics(metricType MetricType, metric string, value float64) *Metrics {
	m := &Metrics{
		ID:    metric,
		MType: string(metricType),
	}

	if metricType == TypeGauge {
		m.Value = &value
	} else {
		v := int64(value)
		m.Delta = &v
	}

	return m
}
