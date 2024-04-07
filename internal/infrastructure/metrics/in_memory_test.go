package metrics

import (
	"github.com/stretchr/testify/assert"
	"testing"
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
