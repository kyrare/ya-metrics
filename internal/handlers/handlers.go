package handlers

import (
	"github.com/kyrare/ya-metrics/internal/metrics"
	"github.com/kyrare/ya-metrics/internal/storage"
	"net/http"
	"strconv"
	"strings"
)

func UpdateHandle(w http.ResponseWriter, r *http.Request, s storage.Storage) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 4 {
		http.NotFound(w, r)
		return
	}

	metricType := metrics.MetricType(parts[1])
	metric := parts[2]
	valueStr := parts[3]

	if metricType != metrics.TypeGauge && metricType != metrics.TypeCounter {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	value, err := strconv.ParseFloat(valueStr, 64)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if metricType == metrics.TypeGauge {
		s.UpdateGauge(metric, value)
	}

	if metricType == metrics.TypeCounter {
		s.UpdateCounter(metric, value)
	}
}
