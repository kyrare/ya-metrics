package main

import (
	"flag"
	"fmt"
	"github.com/kyrare/ya-metrics/internal/handlers"
	"github.com/kyrare/ya-metrics/internal/utils"
	"net/http"
)

func main() {
	addr := utils.GetParameter("a", "ADDRESS", "0.0.0.0:8080", "Адрес сервера (по умолчанию 0.0.0.0:8080)")

	flag.Parse()

	fmt.Println(*addr)

	r := handlers.ServerRouter()

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		panic(err)
	}
}
