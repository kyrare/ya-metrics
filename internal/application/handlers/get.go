package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"net/http"
	"strconv"
)

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	metricType := metrics.MetricType(chi.URLParam(r, "metricType"))
	metric := chi.URLParam(r, "metric")

	if metricType != metrics.TypeGauge && metricType != metrics.TypeCounter {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	value, ok := h.storage.Get(metricType, metric)

	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err := w.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))

	if err != nil {
		fmt.Println(err)
	}
}
