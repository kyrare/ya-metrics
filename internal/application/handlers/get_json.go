package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
)

func (h *Handler) GetJSON(w http.ResponseWriter, r *http.Request) {
	var request metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	metricType := metrics.MetricType(request.MType)

	if metricType != metrics.TypeGauge && metricType != metrics.TypeCounter {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	value, ok := h.storage.Get(metricType, request.ID)

	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	response, err := metrics.NewMetrics(metricType, request.ID, value)

	if err != nil {
		h.logger.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	responseJSON, err := json.Marshal(response)

	if err != nil {
		h.logger.Error("Не удалось сконвертировать response в json", response)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)

	if err != nil {
		h.logger.Error(err)
	}
}
