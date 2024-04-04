package main

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/metrics"
	"github.com/kyrare/ya-metrics/internal/storage/agent"
	"math/rand"
	"net/http"
	"time"
)

func main() {
	storage := agent.MemStorage{
		Values: make(map[string]float64),
	}

	lastSendTime := time.Now()

	for {
		saveRuntimes(&storage)
		storage.Increment("PollCount")
		storage.Set("RandomValue", rand.Float64())

		if time.Since(lastSendTime) >= (10 * time.Second) {
			err := sendStorage(storage)

			if err != nil {
				panic(err)
			}

			lastSendTime = time.Now()
		}

		time.Sleep(2 * time.Second)
	}
}

func saveRuntimes(storage agent.Storage) {
	values := metrics.GetRuntimes()

	for metric, value := range values {
		storage.Set(metric, value)
	}
}

func sendStorage(storage agent.MemStorage) error {
	for metric, value := range storage.Values {
		err := sendMetric(metrics.TypeGauge, metric, value)

		if err != nil {
			return err
		}
	}

	return nil
}

func sendMetric(metricType metrics.MetricType, metric string, value float64) error {
	response, err := http.Post(
		fmt.Sprintf("http://localhost:8080/update/%s/%s/%v", metricType, metric, value),
		"text/plain",
		nil,
	)

	if err != nil {
		return err
	}

	err = response.Body.Close()

	if err != nil {
		return err
	}

	return nil
}
