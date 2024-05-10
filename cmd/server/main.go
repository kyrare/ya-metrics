package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kyrare/ya-metrics/internal/infrastructure/connection"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"github.com/kyrare/ya-metrics/internal/service/server"
	"go.uber.org/zap"
)

func main() {
	config, err := server.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	logger, err := createLogger(config)
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar := *logger.Sugar()

	storage, err := metrics.NewMemStorage(config.FileStoragePath, sugar)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := storage.StoreAndClose()
		if err != nil {
			log.Fatal(err)
		}
	}()

	db, err := connection.New(config.DatabaseDsn, sugar)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		err := storage.StoreAndClose()
		if err != nil {
			sugar.Error(err)
		}
		os.Exit(0)
	}()

	service := server.NewServer(config, storage, db, sugar)

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
