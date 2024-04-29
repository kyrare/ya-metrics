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
	logger  zap.SugaredLogger
}

func NewHandler(storage metrics.Storage, logger zap.SugaredLogger) *Handler {
	return &Handler{storage: storage, logger: logger}
}

func ServerRouter(logger zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()

	s := metrics.NewMemStorage()
	h := NewHandler(s, logger)

	r.Use(func(handler http.Handler) http.Handler {
		return middlewares.WithLogging(handler, logger)
	})

	r.Get("/", h.Home)
	r.Get("/value/{metricType}/{metric}", h.Get)
	r.Post("/update/{metricType}/{metric}/{value}", h.Update)

	r.Post("/value/", h.GetJson)
	r.Post("/update/", h.UpdateJson)
	return r
}
