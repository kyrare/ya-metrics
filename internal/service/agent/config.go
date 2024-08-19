package agent

import (
	"flag"
	"strconv"
	"time"

	"github.com/kyrare/ya-metrics/internal/domain/utils"
)

type Config struct {
	Address        string
	ReportInterval time.Duration
	PollInterval   time.Duration
	AppEnv         string
	AddKey         bool
	RateLimit      uint64
	CryptoKey      string
}

// LoadConfig загружает конфиг для агента
func LoadConfig() (Config, error) {
	addr := utils.GetParameter("a", "ADDRESS", "0.0.0.0:8080", "Адрес сервера (по умолчанию 0.0.0.0:8080)")
	reportIntervalStr := utils.GetParameter("r", "REPORT_INTERVAL", "10", "Частота отправки метрик на сервер (по умолчанию 10 секунд)")
	pollIntervalStr := utils.GetParameter("p", "POLL_INTERVAL", "2", "Частота опроса метрик (по умолчанию 2 секунды)")
	appEnv := utils.GetParameter("env", "APP_ENV", "development", "Режим работы, production|development (по умолчанию development)")
	key := utils.GetParameter("k", "KEY", "", "Добавлять заголовок с хешом")
	rateLimitStr := utils.GetParameter("l", "RATE_LIMIT", "1", "Количество одновременно исходящих запросов на сервер")
	cryptoKey := utils.GetParameter("crypto-key", "CRYPTO_KEY", "", "Путь до файла с публичным ключом")

	flag.Parse()

	reportInterval, err := strconv.Atoi(*reportIntervalStr)
	if err != nil {
		return Config{}, err
	}

	pollInterval, err := strconv.Atoi(*pollIntervalStr)
	if err != nil {
		return Config{}, err
	}

	rateLimit, err := strconv.ParseInt(*rateLimitStr, 10, strconv.IntSize)
	if err != nil {
		return Config{}, err
	}

	return Config{
		Address:        *addr,
		ReportInterval: time.Duration(reportInterval),
		PollInterval:   time.Duration(pollInterval),
		AppEnv:         *appEnv,
		AddKey:         *key != "",
		RateLimit:      uint64(rateLimit),
		CryptoKey:      *cryptoKey,
	}, nil
}
