// Пакет agent для сбора метрик машины и отправки собранных метрик на сервер

package agent

import (
	"math/rand"
	"sync"
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

// NewAgent Конструктор для агента
func NewAgent(config Config, storage agent.Storage, client client.Client, logger zap.SugaredLogger) *Agent {
	return &Agent{
		Config:  config,
		Storage: storage,
		Client:  client,
		Logger:  logger,
	}
}

// Run Запускает агента
func (a *Agent) Run() {
	var wg sync.WaitGroup

	wg.Add(3)

	go func() {
		ticker := time.NewTicker(a.Config.PollInterval * time.Second)

		for range ticker.C {
			a.saveRuntimes()
			a.Storage.IncrementCounter()
			a.Storage.Set("RandomValue", rand.Float64())
		}

		wg.Done()
	}()

	// не знаю зачем тут отдельная горутина, но коль просят по заданию...
	go func() {
		ticker := time.NewTicker(a.Config.PollInterval * time.Second)

		for range ticker.C {
			a.savePstil()
		}

		wg.Done()
	}()

	go func() {
		ticker := time.NewTicker(a.Config.ReportInterval * time.Second)

		for range ticker.C {
			err := a.sendStorage()

			if err != nil {
				a.Logger.Error(err)
			}

			a.Storage.ResetCounter()
		}

		wg.Done()
	}()

	wg.Wait()
}

func (a *Agent) saveRuntimes() {
	values := metrics.GetRuntimes()

	for metric, value := range values {
		a.Storage.Set(metric, value)
	}
}

func (a *Agent) savePstil() {
	values := metrics.GetPstil()

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

	a.Client.Send(data)

	return nil
}
