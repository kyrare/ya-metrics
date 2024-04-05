package main

import (
	"github.com/kyrare/ya-metrics/internal/handlers"
	"net/http"
)

func main() {
	r := handlers.ServerRouter()

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err)
	}
}
