package agent

import (
	"math/rand"
	"time"

	"github.com/kyrare/ya-metrics/internal/application/client"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"github.com/kyrare/ya-metrics/internal/infrastructure/agent"
	"go.uber.org/zap"
)

type Agent struct {
	Config  Config
	Storage agent.Storage
	Client  client.Client
	Logger  zap.SugaredLogger
}

func (a *Agent) Run() {
	lastSendTime := time.Now()

	for {
		a.saveRuntimes()
		a.Storage.IncrementCounter()
		a.Storage.Set("RandomValue", rand.Float64())

		if time.Since(lastSendTime) >= (a.Config.ReportInterval * time.Second) {
			err := a.sendStorage()

			if err != nil {
				a.Logger.Error(err)
			}

			a.Storage.ResetCounter()
			lastSendTime = time.Now()
		}

		time.Sleep(a.Config.PollInterval * time.Second)
	}
}

func (a *Agent) saveRuntimes() {
	values := metrics.GetRuntimes()

	for metric, value := range values {
		a.Storage.Set(metric, value)
	}
}

func (a *Agent) sendStorage() error {
	var data []metrics.Metrics

	for metric, value := range a.Storage.All() {
		m, err := metrics.NewMetrics(metrics.TypeGauge, metric, value)
		if err != nil {
			return err
		}

		data = append(data, *m)
	}

	m, err := metrics.NewMetrics(metrics.TypeCounter, "PollCount", float64(a.Storage.GetCounter()))
	if err != nil {
		return err
	}

	data = append(data, *m)

	return a.Client.Send(data)
}

func NewAgent(config Config, storage agent.Storage, client client.Client, logger zap.SugaredLogger) *Agent {
	return &Agent{
		Config:  config,
		Storage: storage,
		Client:  client,
		Logger:  logger,
	}
}
