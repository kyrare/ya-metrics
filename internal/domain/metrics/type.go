package metrics

import (
	"fmt"

	pb "github.com/kyrare/ya-metrics/internal/infrastructure/proto"
)

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
		return nil, fmt.Errorf("неизвестный тип метрики %s", metricType)
	}

	return m, nil
}

func FromGRPC(m Metrics) *pb.Metric {
	var t pb.Metric_MType
	if m.MType == string(TypeGauge) {
		t = pb.Metric_GAUGE
	}
	if m.MType == string(TypeCounter) {
		t = pb.Metric_COUNTER
	}

	var delta int64
	if m.Delta != nil {
		delta = *m.Delta
	}

	var value float64
	if m.Value != nil {
		value = *m.Value
	}

	return &pb.Metric{
		Id:    m.ID,
		MType: t,
		Delta: delta,
		Value: value,
	}
}

func ToGRPC(m *pb.Metric) Metrics {
	var t MetricType

	if m.MType == pb.Metric_GAUGE {
		t = TypeGauge
	}
	if m.MType == pb.Metric_COUNTER {
		t = TypeCounter
	}

	return Metrics{
		ID:    m.Id,
		MType: string(t),
		Delta: &m.Delta,
		Value: &m.Value,
	}
}
