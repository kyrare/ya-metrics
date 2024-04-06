package main

import (
	"flag"
	"github.com/kyrare/ya-metrics/internal/handlers"
	"net/http"
)

func main() {
	addr := flag.String("a", "0.0.0.0:8080", "Адрес сервера (по умолчанию 0.0.0.0:8080)")
	flag.Parse()

	//fmt.Println(*addr)

	r := handlers.ServerRouter()

	err := http.ListenAndServe(*addr, r)
	if err != nil {
		panic(err)
	}
}
