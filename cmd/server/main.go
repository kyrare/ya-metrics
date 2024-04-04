package main

import (
	"github.com/kyrare/ya-metrics/internal/handlers"
	"github.com/kyrare/ya-metrics/internal/middlewares"
	"github.com/kyrare/ya-metrics/internal/storage"
	"net/http"
)

func main() {
	var s storage.MemStorage

	mux := http.NewServeMux()

	mux.Handle("/update/", middlewares.PostMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.UpdateHandle(w, r, &s)
	})))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
