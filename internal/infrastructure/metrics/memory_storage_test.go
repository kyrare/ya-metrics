package metrics

import (
	"testing"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		Gauges   map[string]float64
		Counters map[string]float64
	}
	type args struct {
		metric string
		value  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]float64
	}{
		{
			name: "add value",
			fields: fields{
				Gauges:   map[string]float64{},
				Counters: map[string]float64{},
			},
			args: args{
				metric: "foo",
				value:  1,
			},
			want: map[string]float64{
				"foo": 1,
			},
		},
		{
			name: "replace value",
			fields: fields{
				Gauges: map[string]float64{
					"foo": 1,
				},
				Counters: map[string]float64{},
			},
			args: args{
				metric: "foo",
				value:  10,
			},
			want: map[string]float64{
				"foo": 10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			storage.UpdateGauge(tt.args.metric, tt.args.value)
			assert.Equal(t, tt.want, storage.Gauges)
		})
	}
}

func TestMemStorage_UpdateCounter(t *testing.T) {
	type fields struct {
		Gauges   map[string]float64
		Counters map[string]float64
	}
	type args struct {
		metric string
		value  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]float64
	}{
		{
			name: "add value",
			fields: fields{
				Gauges:   map[string]float64{},
				Counters: map[string]float64{},
			},
			args: args{
				metric: "foo",
				value:  1,
			},
			want: map[string]float64{
				"foo": 1,
			},
		},
		{
			name: "add second value",
			fields: fields{
				Gauges: map[string]float64{},
				Counters: map[string]float64{
					"foo": 1,
				},
			},
			args: args{
				metric: "foo",
				value:  10,
			},
			want: map[string]float64{
				"foo": 11,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			storage := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}
			storage.UpdateCounter(tt.args.metric, tt.args.value)
			assert.Equal(t, tt.want, storage.Counters)
		})
	}
}

func TestMemStorage_Updates(t *testing.T) {
	type fields struct {
		Gauges   map[string]float64
		Counters map[string]float64
	}

	type args struct {
		values []metrics.Metrics
	}

	value := float64(1.1)
	delta := int64(1)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "Test Updates",
			fields: fields{
				Gauges:   map[string]float64{},
				Counters: map[string]float64{},
			},
			args: args{
				values: []metrics.Metrics{
					{
						ID:    "TEST_GAUGE",
						MType: string(metrics.TypeGauge),
						Value: &value,
					},
					{
						ID:    "TEST_COUNTER",
						MType: string(metrics.TypeCounter),
						Delta: &delta,
					},
				},
			},
			want: fields{
				Gauges: map[string]float64{
					"TEST_GAUGE": 1.1,
				},
				Counters: map[string]float64{
					"TEST_COUNTER": 1,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}

			s.Updates(tt.args.values)

			assert.Equal(t, tt.want, fields{
				Gauges:   s.GetGauges(),
				Counters: s.GetCounters(),
			})
		})
	}
}

func TestMemStorage_GetValue(t *testing.T) {
	type fields struct {
		Gauges   map[string]float64
		Counters map[string]float64
	}
	type args struct {
		metricType metrics.MetricType
		metric     string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   float64
		want1  bool
	}{
		{
			name: "Get Gauges Value",
			fields: fields{
				Gauges: map[string]float64{
					"TEST_1": 1.1,
				},
				Counters: map[string]float64{},
			},
			args: args{
				metricType: metrics.TypeGauge,
				metric:     "TEST_1",
			},
			want:  1.1,
			want1: true,
		},
		{
			name: "Get Counters Value",
			fields: fields{
				Gauges: map[string]float64{},
				Counters: map[string]float64{
					"TEST_2": 2.2,
				},
			},
			args: args{
				metricType: metrics.TypeCounter,
				metric:     "TEST_2",
			},
			want:  2.2,
			want1: true,
		},
		{
			name: "Empty Gauges Value",
			fields: fields{
				Gauges:   map[string]float64{},
				Counters: map[string]float64{},
			},
			args: args{
				metricType: metrics.TypeCounter,
				metric:     "TEST_3",
			},
			want:  0,
			want1: false,
		},
		{
			name: "Empty Counters Value",
			fields: fields{
				Gauges:   map[string]float64{},
				Counters: map[string]float64{},
			},
			args: args{
				metricType: metrics.TypeCounter,
				metric:     "TEST_4",
			},
			want:  0,
			want1: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &MemStorage{
				Gauges:   tt.fields.Gauges,
				Counters: tt.fields.Counters,
			}

			got, got1 := s.GetValue(tt.args.metricType, tt.args.metric)

			assert.Equalf(t, tt.want, got, "GetValue(%v, %v)", tt.args.metricType, tt.args.metric)
			assert.Equalf(t, tt.want1, got1, "GetValue(%v, %v)", tt.args.metricType, tt.args.metric)
		})
	}
}
