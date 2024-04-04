package main

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/agent_storage"
	"github.com/kyrare/ya-metrics/internal/metrics"
	"math/rand/v2"
	"net/http"
	"time"
)

func main() {
	storage := agent_storage.MemStorage{
		Values: make(map[string]float64),
	}

	lastSendTime := time.Now()

	for {
		saveRuntimes(&storage)
		storage.Increment("PollCount")
		storage.Set("RandomValue", rand.Float64())

		if time.Now().Sub(lastSendTime) >= (10 * time.Second) {
			sendStorage(storage)
			lastSendTime = time.Now()
		}

		time.Sleep(2 * time.Second)
	}
}

func saveRuntimes(storage agent_storage.Storage) {
	values := metrics.GetRuntimes()

	for metric, value := range values {
		storage.Set(metric, value)
	}
}

func sendStorage(storage agent_storage.MemStorage) {
	for metric, value := range storage.Values {
		sendMetric(metrics.TypeGauge, metric, value)
	}
}

func sendMetric(metricType metrics.MetricType, metric string, value float64) {
	_, err := http.Post(
		fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metricType, metric, value),
		"text/plain",
		nil,
	)

	if err != nil {
		panic(err)
	}
}
