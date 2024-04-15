package main

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/service/server"
	"log"
)

func main() {
	config, err := server.LoadConfig()

	if err != nil {
		log.Fatal(err)
	}

	service := server.NewServer(config)

	err = service.Run()

	if err != nil {
		fmt.Println(err)
	}
}
