package main

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/service/server"
	"go.uber.org/zap"
	"log"
)

func main() {
	config, err := server.LoadConfig()

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

	service := server.NewServer(config, sugar)

	err = service.Run()

	if err != nil {
		fmt.Println(err)
	}
}
