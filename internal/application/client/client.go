package client

import "github.com/kyrare/ya-metrics/internal/domain/metrics"

type Client interface {
	Send(data []metrics.Metrics)
	Close() error
}
