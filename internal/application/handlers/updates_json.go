package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/kyrare/ya-metrics/internal/domain/metrics"
)

// UpdatesJSON обновление нескольких метрик в формате JSON
func (h *Handler) UpdatesJSON(w http.ResponseWriter, r *http.Request) {
	if !h.checkRequestKey(r) {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	var request []metrics.Metrics

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := h.storage.Updates(request); err != nil {
		h.logger.Error("Не удалось обновить записи", err, request)
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)

	if h.storeStorageOnHit {
		err := h.storage.Store()
		if err != nil {
			h.logger.Error(err)
		}
	}
}
