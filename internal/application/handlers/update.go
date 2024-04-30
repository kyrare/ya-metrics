package handlers

import (
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
)

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	metricType := metrics.MetricType(chi.URLParam(r, "metricType"))

	if metricType != metrics.TypeGauge && metricType != metrics.TypeCounter {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	metric := chi.URLParam(r, "metric")
	value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if metricType == metrics.TypeGauge {
		h.storage.UpdateGauge(metric, value)
	}

	if metricType == metrics.TypeCounter {
		h.storage.UpdateCounter(metric, value)
	}

	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	if h.storeStorageOnHit {
		err := h.storage.Store()
		if err != nil {
			h.logger.Error(err)
		}
	}
}
