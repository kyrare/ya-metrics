package main

import (
	"fmt"
	"log"

	"github.com/kyrare/ya-metrics/internal/application/client"
	"github.com/kyrare/ya-metrics/internal/domain/utils"
	agentStorage "github.com/kyrare/ya-metrics/internal/infrastructure/agent"
	"github.com/kyrare/ya-metrics/internal/service/agent"
	"go.uber.org/zap"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	utils.PrintBuildData(buildVersion, buildDate, buildCommit)

	config, err := agent.LoadConfig()

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

	s := agentStorage.NewMemStorage()

	var cl client.Client

	if config.UseGRPC {
		sugar.Info("Agent use gRPC")
		cl, err = client.NewGRPCClient(config.Address, sugar)
	} else {
		sugar.Info("Agent use HTTP")
		cl, err = client.NewHTTPClient(config.Address, config.AddKey, config.RateLimit, config.CryptoKey, sugar)
	}
	if err != nil {
		log.Fatal(err)
	}
	defer cl.Close()

	service := agent.NewAgent(config, s, cl, sugar)

	service.Run()
}

func createLogger(c agent.Config) (*zap.Logger, error) {
	switch c.AppEnv {
	case "development":
		return zap.NewDevelopment()
	case "production":
		return zap.NewProduction()
	default:
		return nil, fmt.Errorf("неизвестный APP_ENV %s", c.AppEnv)
	}
}
