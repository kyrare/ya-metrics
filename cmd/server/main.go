package main

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/service/server"
)

func main() {
	config, err := server.LoadConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	service := server.NewServer(config)

	err = service.Run()

	if err != nil {
		fmt.Println(err)
	}
}
