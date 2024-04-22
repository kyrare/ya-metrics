package main

import (
	"github.com/kyrare/ya-metrics/internal/application/client"
	agentStorage "github.com/kyrare/ya-metrics/internal/infrastructure/agent"
	"github.com/kyrare/ya-metrics/internal/service/agent"
	"go.uber.org/zap"
	"log"
)

func main() {
	c, err := agent.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	// делаем регистратор SugaredLogger
	sugar := *logger.Sugar()

	s := agentStorage.NewMemeStorage()
	cl := client.NewClient(c.Address, sugar)

	service := agent.NewAgent(c, s, *cl)

	service.Run()
}
