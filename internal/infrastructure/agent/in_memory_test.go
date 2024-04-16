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
				values: tt.fields.values,
			}

			m.Set(tt.args.metric, tt.args.value)

			assert.Contains(t, m.values, tt.args.metric)
			assert.Equal(t, m.values[tt.args.metric], tt.args.value)
		})
	}
}

func TestMemStorage_increment(t *testing.T) {
	type fields struct {
		counter int
	}

	tests := []struct {
		name   string
		fields fields
		actual int
	}{
		{
			name: "Test 0 -> 1",
			fields: fields{
				counter: 0,
			},
			actual: 1,
		},
		{
			name: "Test 1000 -> 1001",
			fields: fields{
				counter: 1000,
			},
			actual: 1001,
		},
		{
			name: "Test -1 -> 0",
			fields: fields{
				counter: -1,
			},
			actual: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				counter: tt.fields.counter,
			}

			m.IncrementCounter()
			assert.Equal(t, m.counter, tt.actual)
		})
	}
}
