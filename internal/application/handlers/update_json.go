package handlers

import (
	"encoding/json"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"net/http"
)

func (h *Handler) UpdateJSON(w http.ResponseWriter, r *http.Request) {
	var request metrics.Metrics

	h.logger.Infof("Test1: %v", request)

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	h.logger.Infof("Test2: %v", request)

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

	responseData, err := metrics.NewMetrics(metricType, request.ID, value)

	if err != nil {
		h.logger.Error(err)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	responseJSON, err := json.Marshal(responseData)

	if err != nil {
		h.logger.Error("Не удалось сконвертировать responseData в json", responseData)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJSON)

	if err != nil {
		h.logger.Error(err)
	}

	if h.storeStorageOnHit {
		err := h.storage.Store()
		if err != nil {
			h.logger.Error(err)
		}
	}
}
