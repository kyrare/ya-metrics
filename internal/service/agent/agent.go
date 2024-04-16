package agent

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/application/client"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"github.com/kyrare/ya-metrics/internal/infrastructure/agent"
	"math/rand"
	"time"
)

type Agent struct {
	Config  Config
	Storage agent.Storage
	Client  client.Client
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
				fmt.Println(err)
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
	for metric, value := range a.Storage.All() {
		err := a.Client.Send(metrics.TypeGauge, metric, value)

		if err != nil {
			return err
		}
	}

	return a.Client.Send(metrics.TypeCounter, "PollCount", float64(a.Storage.GetCounter()))
}

func NewAgent(config Config, storage agent.Storage, client client.Client) *Agent {
	return &Agent{
		Config:  config,
		Storage: storage,
		Client:  client,
	}
}
