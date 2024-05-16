package handlers

import (
	"database/sql"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/kyrare/ya-metrics/internal/application/middlewares"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"go.uber.org/zap"
)

type Handler struct {
	storage           metrics.Storage
	storeStorageOnHit bool
	DB                *sql.DB
	logger            zap.SugaredLogger
}

func NewHandler(storage metrics.Storage, DB *sql.DB, storeDataOnHit bool, logger zap.SugaredLogger) *Handler {
	return &Handler{
		storage:           storage,
		storeStorageOnHit: storeDataOnHit,
		DB:                DB,
		logger:            logger,
	}
}

func ServerRouter(storage metrics.Storage, DB *sql.DB, storeDataOnHit bool, logger zap.SugaredLogger) chi.Router {
	r := chi.NewRouter()

	h := NewHandler(storage, DB, storeDataOnHit, logger)

	r.Use(func(handler http.Handler) http.Handler {
		return middlewares.WithLogging(handler, logger)
	})
	r.Use(middlewares.Compress)
	r.Use(middlewares.Decompress)

	r.Get("/", h.Home)
	r.Get("/ping", h.Ping)
	r.Get("/value/{metricType}/{metric}", h.Get)
	r.Post("/update/{metricType}/{metric}/{value}", h.Update)

	r.Post("/value/", h.GetJSON)
	r.Post("/update/", h.UpdateJSON)

	r.Post("/updates/", h.UpdatesJSON)
	return r
}
