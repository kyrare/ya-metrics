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
	buildVersion string
	buildDate    string
	buildCommit  string
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
	cl := client.NewClient(config.Address, config.AddKey, config.RateLimit, sugar)

	service := agent.NewAgent(config, s, *cl, sugar)

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
