package handlers

import (
	"encoding/json"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"net/http"
)

func (h *Handler) UpdateJson(w http.ResponseWriter, r *http.Request) {
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

	if metricType == metrics.TypeGauge && request.Value == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if metricType == metrics.TypeCounter && request.Delta == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if metricType == metrics.TypeGauge {
		h.storage.UpdateGauge(request.ID, *request.Value)
	} else {
		h.storage.UpdateCounter(request.ID, float64(*request.Delta))
	}

	value, _ := h.storage.Get(metricType, request.ID)

	responseData := metrics.NewMetrics(metricType, request.ID, value)

	responseJson, err := json.Marshal(responseData)

	if err != nil {
		h.logger.Error("Не удалось сконвертировать responseData в json", responseData)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJson)

	if err != nil {
		h.logger.Error(err)
	}
}
