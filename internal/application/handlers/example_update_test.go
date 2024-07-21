package handlers

import (
	"github.com/go-chi/chi/v5"
	"github.com/kyrare/ya-metrics/internal/infrastructure/connection"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"go.uber.org/zap"
)

func ExampleHandler_Update() {
	logger := zap.New(nil)
	sugar := *logger.Sugar()
	storage, err := metrics.NewMemStorage("", sugar)
	if err != nil {
		sugar.Error(err)
		return
	}

	db, err := connection.New("", sugar)
	if err != nil {
		sugar.Error(err)
		return
	}

	r := chi.NewRouter()
	h := NewHandler(storage, db, false, false, sugar)

	r.Get("/update/{metricType}/{metric}/{value}", h.Update)
}
