package agent

import (
	"encoding/json"
	"flag"
	"io"
	"os"
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
	UseGRPC        bool
}

type configFile struct {
	Address        string `json:"address,omitempty"`
	ReportInterval string `json:"store_interval,omitempty"`
	PollInterval   string `json:"poll_interval,omitempty"`
	Key            string `json:"key,omitempty"`
	AppEnv         string `json:"app_env,omitempty"`
	RateLimit      string `json:"rate_limit,omitempty"`
	CryptoKey      string `json:"crypto_key,omitempty"`
	UseGRPC        string `json:"use_grpc,omitempty"`
}

// LoadConfig загружает конфиг для агента
func LoadConfig() (Config, error) {
	configFilePath := utils.GetParameter("c", "CONFIG", "", "", "Путь до файла с конфигурацией")

	cf, err := loadConfigFile(*configFilePath)
	if err != nil {
		return Config{}, err
	}

	addr := utils.GetParameter("a", "ADDRESS", cf.Address, "0.0.0.0:8080", "Адрес сервера (по умолчанию 0.0.0.0:8080)")
	reportIntervalStr := utils.GetParameter("r", "REPORT_INTERVAL", cf.ReportInterval, "10", "Частота отправки метрик на сервер (по умолчанию 10 секунд)")
	pollIntervalStr := utils.GetParameter("p", "POLL_INTERVAL", cf.PollInterval, "2", "Частота опроса метрик (по умолчанию 2 секунды)")
	appEnv := utils.GetParameter("env", "APP_ENV", cf.AppEnv, "development", "Режим работы, production|development (по умолчанию development)")
	key := utils.GetParameter("k", "KEY", cf.Key, "", "Добавлять заголовок с хешом")
	rateLimitStr := utils.GetParameter("l", "RATE_LIMIT", cf.RateLimit, "1", "Количество одновременно исходящих запросов на сервер")
	cryptoKey := utils.GetParameter("crypto-key", "CRYPTO_KEY", cf.CryptoKey, "", "Путь до файла с публичным ключом")
	useGRPC := utils.GetParameter("grpc", "USE_GRPC", cf.UseGRPC, "", "Использовать протокол gRPC")

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
		UseGRPC:        *useGRPC != "",
	}, nil
}

func loadConfigFile(path string) (configFile, error) {
	if path == "" {
		return configFile{}, nil
	}

	file, err := os.Open(path)
	if err != nil {
		return configFile{}, err
	}
	defer file.Close()

	data, _ := io.ReadAll(file)

	var cf configFile
	err = json.Unmarshal(data, &cf)
	if err != nil {
		return configFile{}, err
	}

	return cf, err
}
