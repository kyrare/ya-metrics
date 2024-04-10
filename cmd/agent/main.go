package main

import (
	"fmt"
	"github.com/kyrare/ya-metrics/internal/application/client"
	agentStorage "github.com/kyrare/ya-metrics/internal/infrastructure/agent"
	"github.com/kyrare/ya-metrics/internal/service/agent"
)

func main() {
	c, err := agent.LoadConfig()

	if err != nil {
		fmt.Println(err)
		return
	}

	s := agentStorage.NewMemeStorage()
	cl := client.NewClient(c.Address)

	service := agent.NewAgent(c, s, *cl)

	service.Run()
}
