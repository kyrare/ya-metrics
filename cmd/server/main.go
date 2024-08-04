package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/kyrare/ya-metrics/internal/domain/utils"
	"github.com/kyrare/ya-metrics/internal/infrastructure/connection"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"github.com/kyrare/ya-metrics/internal/service/server"
	"go.uber.org/zap"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	utils.PrintBuildData(buildVersion, buildDate, buildCommit)

	ctx := context.Background()

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

	storage, err := createStorage(ctx, config, db, sugar)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := storage.StoreAndClose()
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

func createStorage(ctx context.Context, c server.Config, db *sql.DB, logger zap.SugaredLogger) (metrics.Storage, error) {
	if c.DatabaseDsn == "" {
		return metrics.NewMemStorage(c.FileStoragePath, logger)
	} else {
		return metrics.NewDatabaseStorage(ctx, db, logger)
	}
}
