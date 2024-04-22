package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kyrare/ya-metrics/internal/application/middlewares"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	storage metrics.Storage
}

func NewHandler(storage metrics.Storage) *Handler {
	return &Handler{storage: storage}
}

func ServerRouter(logger zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()

	s := metrics.NewMemStorage()
	h := NewHandler(s)

	r.Use(func(handler http.Handler) http.Handler {
		return middlewares.WithLogging(handler, logger)
	})

	r.Get("/", h.Home)
	r.Get("/value/{metricType}/{metric}", h.Get)
	r.Post("/update/{metricType}/{metric}/{value}", h.Update)

	return r
}
