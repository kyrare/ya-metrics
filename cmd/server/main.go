package main

import (
	"github.com/kyrare/ya-metrics/internal/service/server"
)

func main() {
	config, err := server.LoadConfig()

	if err != nil {
		panic(err)
	}

	service := server.NewServer(config)

	err = service.Run()

	if err != nil {
		panic(err)
	}
}
