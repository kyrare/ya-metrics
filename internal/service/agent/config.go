package agent

import (
	"flag"
	"github.com/kyrare/ya-metrics/internal/domain/utils"
	"strconv"
	"time"
)

type Config struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
}

func LoadConfig() (Config, error) {
	addr := utils.GetParameter("a", "ADDRESS", "0.0.0.0:8080", "Адрес сервера (по умолчанию 0.0.0.0:8080)")
	reportIntervalStr := utils.GetParameter("r", "REPORT_INTERVAL", "10", "Частота отправки метрик на сервер (по умолчанию 10 секунд)")
	pollIntervalStr := utils.GetParameter("p", "POLL_INTERVAL", "2", "Частота опроса метрик (по умолчанию 2 секунды)")

	flag.Parse()

	reportInterval, err := strconv.Atoi(*reportIntervalStr)
	if err != nil {
		return Config{}, err
	}

	pollInterval, err := strconv.Atoi(*pollIntervalStr)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Address:        *addr,
		ReportInterval: time.Duration(reportInterval),
		PollInterval:   time.Duration(pollInterval),
	}, nil
}
