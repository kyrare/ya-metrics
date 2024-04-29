package handlers

import (
	"encoding/json"
	"github.com/kyrare/ya-metrics/internal/domain/metrics"
	"net/http"
	"net/http/httputil"
)

func (h *Handler) UpdateJson(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")

	var request metrics.Metrics

	res, _ := httputil.DumpRequest(r, true)

	h.logger.Info("UpdateJson", string(res))

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var mvalue float64
	if request.Value != nil {
		mvalue = *request.Value
	}

	var delta int64
	if request.Delta != nil {
		delta = *request.Delta
	}

	h.logger.Info("UpdateJson", " ", request.MType, " ", request.ID, " ", mvalue, " ", delta)

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

	w.WriteHeader(http.StatusOK)

	_, err = w.Write(responseJson)

	if err != nil {
		h.logger.Error(err)
	}
}
