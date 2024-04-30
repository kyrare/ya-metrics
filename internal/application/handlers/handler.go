package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kyrare/ya-metrics/internal/application/middlewares"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"go.uber.org/zap"
	"net/http"
)

type Handler struct {
	storage           metrics.Storage
	storeStorageOnHit bool
	logger            zap.SugaredLogger
}

func NewHandler(storage metrics.Storage, storeDataOnHit bool, logger zap.SugaredLogger) *Handler {
	return &Handler{
		storage:           storage,
		storeStorageOnHit: storeDataOnHit,
		logger:            logger,
	}
}

func ServerRouter(storage metrics.Storage, storeDataOnHit bool, logger zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()

	h := NewHandler(storage, storeDataOnHit, logger)

	r.Use(func(handler http.Handler) http.Handler {
		return middlewares.WithLogging(handler, logger)
	})
	r.Use(middlewares.Compress)
	r.Use(middlewares.Decompress)

	r.Get("/", h.Home)
	r.Get("/value/{metricType}/{metric}", h.Get)
	r.Post("/update/{metricType}/{metric}/{value}", h.Update)

	r.Post("/value/", h.GetJSON)
	r.Post("/update/", h.UpdateJSON)
	return r
}
