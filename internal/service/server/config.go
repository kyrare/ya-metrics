package server

import (
	"flag"
	"strconv"
	"time"

	"github.com/kyrare/ya-metrics/internal/domain/utils"
)

type Config struct {
	Address         string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	AppEnv          string
	DatabaseDsn     string
	CheckKey        bool
	CryptoKey       string
}

// LoadConfig загружает конфиг для сервера
func LoadConfig() (Config, error) {
	addr := utils.GetParameter("a", "ADDRESS", "0.0.0.0:8080", "Адрес сервера (по умолчанию 0.0.0.0:8080)")
	storeIntervalStr := utils.GetParameter("i", "STORE_INTERVAL", "300", "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)")
	fileStoragePath := utils.GetParameter("f", "FILE_STORAGE_PATH", "/tmp/metrics-db.json", "Полное имя файла, куда сохраняются текущие значения (по умолчанию /tmp/metrics-db.json, пустое значение отключает функцию записи на диск)")
	restore := utils.GetParameter("r", "RESTORE", "true", "Булево значение (true/false), определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера (по умолчанию true)")
	appEnv := utils.GetParameter("env", "APP_ENV", "development", "Режим работы, production|development (по умолчанию development)")
	databaseDsn := utils.GetParameter("d", "DATABASE_DSN", "", "Строка с адресом подключения к БД")
	key := utils.GetParameter("k", "KEY", "", "Проверять заголовок с хешом")
	cryptoKey := utils.GetParameter("crypto-key", "CRYPTO_KEY", "", "Путь до файла с приватным ключом")

	storeInterval, err := strconv.Atoi(*storeIntervalStr)
	if err != nil {
		return Config{}, err
	}

	flag.Parse()

	return Config{
		Address:         *addr,
		StoreInterval:   time.Duration(storeInterval),
		FileStoragePath: *fileStoragePath,
		Restore:         *restore == "true",
		AppEnv:          *appEnv,
		DatabaseDsn:     *databaseDsn,
		CheckKey:        *key != "",
		CryptoKey:       *cryptoKey,
	}, nil
}
