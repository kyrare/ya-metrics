package server

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
	Address         string
	StoreInterval   time.Duration
	FileStoragePath string
	Restore         bool
	AppEnv          string
	DatabaseDsn     string
	CheckKey        bool
	CryptoKey       string
	TrustedSubnet   string
	UseGRPC         bool
}

type configFile struct {
	Address       string `json:"address,omitempty"`
	Restore       string `json:"restore,omitempty"`
	StoreInterval string `json:"store_interval,omitempty"`
	StoreFile     string `json:"store_file,omitempty"`
	DatabaseDsn   string `json:"database_dsn,omitempty"`
	CryptoKey     string `json:"crypto_key,omitempty"`
	Key           string `json:"key,omitempty"`
	AppEnv        string `json:"app_env,omitempty"`
	TrustedSubnet string `json:"trusted_subnet,omitempty"`
	UseGRPC       string `json:"use_grpc,omitempty"`
}

// LoadConfig загружает конфиг для сервера
func LoadConfig() (Config, error) {
	configFilePath := utils.GetParameter("c", "CONFIG", "", "", "Путь до файла с конфигурацией")

	cf, err := loadConfigFile(*configFilePath)
	if err != nil {
		return Config{}, err
	}

	addr := utils.GetParameter("a", "ADDRESS", cf.Address, "0.0.0.0:8080", "Адрес сервера (по умолчанию 0.0.0.0:8080)")
	storeIntervalStr := utils.GetParameter("i", "STORE_INTERVAL", cf.StoreInterval, "300", "Интервал времени в секундах, по истечении которого текущие показания сервера сохраняются на диск (по умолчанию 300 секунд, значение 0 делает запись синхронной)")
	fileStoragePath := utils.GetParameter("f", "FILE_STORAGE_PATH", cf.StoreFile, "/tmp/metrics-db.json", "Полное имя файла, куда сохраняются текущие значения (по умолчанию /tmp/metrics-db.json, пустое значение отключает функцию записи на диск)")
	restore := utils.GetParameter("r", "RESTORE", cf.Restore, "true", "Булево значение (true/false), определяющее, загружать или нет ранее сохранённые значения из указанного файла при старте сервера (по умолчанию true)")
	appEnv := utils.GetParameter("env", "APP_ENV", cf.AppEnv, "development", "Режим работы, production|development (по умолчанию development)")
	databaseDsn := utils.GetParameter("d", "DATABASE_DSN", cf.DatabaseDsn, "", "Строка с адресом подключения к БД")
	key := utils.GetParameter("k", "KEY", cf.Key, "", "Проверять заголовок с хешом")
	cryptoKey := utils.GetParameter("crypto-key", "CRYPTO_KEY", cf.CryptoKey, "", "Путь до файла с приватным ключом")
	trustedSubnet := utils.GetParameter("t", "TRUSTED_SUBNET", cf.TrustedSubnet, "", "Разрешенные IP адреса (CIDR)")
	useGRPC := utils.GetParameter("grpc", "USE_GRPC", cf.UseGRPC, "", "Использовать протокол gRPC")

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
		TrustedSubnet:   *trustedSubnet,
		UseGRPC:         *useGRPC != "",
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
