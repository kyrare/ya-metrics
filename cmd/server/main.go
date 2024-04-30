package main

import (
	"fmt"
	"log"

	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"github.com/kyrare/ya-metrics/internal/service/server"
	"go.uber.org/zap"
)

func main() {
	config, err := server.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	storage := metrics.NewMemStorage(config.FileStoragePath)

	logger, err := createLogger(config)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar := *logger.Sugar()

	service := server.NewServer(config, storage, sugar)

	err = service.Run()

	if err != nil {
		fmt.Println(err)
	}
}

func createLogger(c server.Config) (*zap.Logger, error) {
	switch c.AppEnv {
	case "development":
		return zap.NewDevelopment()
	case "production":
		return zap.NewProduction()
	default:
		return nil, fmt.Errorf("неизвестный APP_ENV %s", c.AppEnv)
	}
}
