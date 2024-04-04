package agent

import (
	"testing"
)
import "github.com/stretchr/testify/assert"

func TestMemStorage_set(t *testing.T) {
	type fields struct {
		values map[string]float64
	}
	type args struct {
		metric string
		value  float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "Test Set",
			fields: fields{
				values: make(map[string]float64),
			},
			args: args{
				metric: "foo",
				value:  10,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Values: tt.fields.values,
			}

			m.Set(tt.args.metric, tt.args.value)

			assert.Contains(t, m.Values, tt.args.metric)
			assert.Equal(t, m.Values[tt.args.metric], tt.args.value)
		})
	}
}

func TestMemStorage_increment(t *testing.T) {
	type fields struct {
		values map[string]float64
	}

	tests := []struct {
		name   string
		fields fields
		metric string
		actual float64
	}{
		{
			name: "Test nil -> 1",
			fields: fields{
				values: make(map[string]float64),
			},
			metric: "foo",
			actual: 1,
		},
		{
			name: "Test 0 -> 1",
			fields: fields{
				values: map[string]float64{
					"foo": 0,
				},
			},
			metric: "foo",
			actual: 1,
		},
		{
			name: "Test 1000 -> 1001",
			fields: fields{
				values: map[string]float64{
					"foo": 1000,
				},
			},
			metric: "foo",
			actual: 1001,
		},
		{
			name: "Test -1 -> 0",
			fields: fields{
				values: map[string]float64{
					"foo": -1,
				},
			},
			metric: "foo",
			actual: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				Values: tt.fields.values,
			}

			m.Increment(tt.metric)
			assert.Contains(t, m.Values, tt.metric)
			assert.Equal(t, m.Values[tt.metric], tt.actual)
		})
	}
}
