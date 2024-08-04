package metrics

import (
	"context"
	"testing"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"github.com/kyrare/ya-metrics/internal/infrastructure/connection"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestDatabaseStorage_Updates(t *testing.T) {
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
			ctx := context.Background()
			logger := *zap.New(nil).Sugar()

			db, err := connection.New("postgres://postgres:postgres@localhost:5433/praktikum_test?sslmode=disable", logger)
			assert.NoError(t, err)

			s, err := NewDatabaseStorage(ctx, db, logger)
			assert.NoError(t, err)

			s.Updates(tt.args.values)

			assert.Equal(t, tt.want, fields{
				Gauges:   s.GetGauges(),
				Counters: s.GetCounters(),
			})
		})
	}
}
