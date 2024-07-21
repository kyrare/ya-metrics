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

func TestMemStorage_IncrementCounter(t *testing.T) {
	type fields struct {
		counter int
	}

	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Test 0 -> 1",
			fields: fields{
				counter: 0,
			},
			want: 1,
		},
		{
			name: "Test 1000 -> 1001",
			fields: fields{
				counter: 1000,
			},
			want: 1001,
		},
		{
			name: "Test -1 -> 0",
			fields: fields{
				counter: -1,
			},
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				counter: tt.fields.counter,
			}

			m.IncrementCounter()
			assert.Equal(t, tt.want, m.GetCounter())
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
	type fields struct {
		counter int
	}

	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "Test 0",
			fields: fields{
				counter: 0,
			},
			want: 0,
		},
		{
			name: "Test 1",
			fields: fields{
				counter: 1,
			},
			want: 1,
		},
		{
			name: "Test 1000",
			fields: fields{
				counter: 1000,
			},
			want: 1000,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &MemStorage{
				counter: tt.fields.counter,
			}

			assert.Equal(t, tt.want, m.GetCounter())
		})
	}
}

func TestMemStorage_ResetCounter(t *testing.T) {
	t.Run("Reset counter", func(t *testing.T) {
		m := &MemStorage{
			counter: 0,
		}

		m.IncrementCounter()
		m.IncrementCounter()
		m.IncrementCounter()

		assert.Equal(t, 3, m.GetCounter())

		m.ResetCounter()

		assert.Equal(t, 0, m.GetCounter())
	})
}
