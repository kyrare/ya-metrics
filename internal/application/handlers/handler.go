package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
)

type Handler struct {
	storage metrics.Storage
}

func NewHandler(storage metrics.Storage) *Handler {
	return &Handler{storage: storage}
}

func ServerRouter() chi.Router {
	r := chi.NewRouter()

	s := metrics.NewMemStorage()
	h := NewHandler(s)

	r.Get("/", h.Home)
	r.Get("/value/{metricType}/{metric}", h.Get)
	r.Post("/update/{metricType}/{metric}/{value}", h.Update)

	return r
}
