package main

import (
	"flag"
	"fmt"
	"github.com/kyrare/ya-metrics/internal/client"
	"github.com/kyrare/ya-metrics/internal/metrics"
	"github.com/kyrare/ya-metrics/internal/storage/agent"
	"github.com/kyrare/ya-metrics/internal/utils"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	addr := utils.GetParameter("a", "ADDRESS", "0.0.0.0:8080", "Адрес сервера (по умолчанию 0.0.0.0:8080)")
	reportIntervalStr := utils.GetParameter("r", "REPORT_INTERVAL", "10", "Частота отправки метрик на сервер (по умолчанию 10 секунд)")
	pollIntervalStr := utils.GetParameter("p", "POLL_INTERVAL", "2", "Частота опроса метрик (по умолчанию 2 секунды)")

	flag.Parse()

	reportInterval, err := strconv.Atoi(*reportIntervalStr)
	if err != nil {
		panic(err)
	}

	pollInterval, err := strconv.Atoi(*pollIntervalStr)
	if err != nil {
		panic(err)
	}

	storage := agent.NewMemeStorage()
	c := client.NewClient(*addr)

	lastSendTime := time.Now()

	for {
		saveRuntimes(storage)
		storage.Increment("PollCount")
		storage.Set("RandomValue", rand.Float64())

		if time.Since(lastSendTime) >= (time.Duration(reportInterval) * time.Second) {
			err := sendStorage(c, storage)

			if err != nil {
				fmt.Println(err)
			}

			lastSendTime = time.Now()
		}

		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}

func saveRuntimes(storage agent.Storage) {
	values := metrics.GetRuntimes()

	for metric, value := range values {
		storage.Set(metric, value)
	}
}

func sendStorage(c *client.Client, storage agent.Storage) error {
	for metric, value := range storage.All() {
		err := c.Send(metrics.TypeGauge, metric, value)

		if err != nil {
			return err
		}
	}

	return nil
}
