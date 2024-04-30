package main

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"github.com/kyrare/ya-metrics/internal/service/server"
	"go.uber.org/zap"
	"log"
)

func main() {
	config, err := server.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	storage := metrics.NewMemStorage(config.FileStoragePath)

	logger, err := zap.NewDevelopment()
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
