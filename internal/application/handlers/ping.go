package handlers

import (
	"net/http"
)

// Ping эндпоинт для проверки работы сайта
func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.DB.Ping()

	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}

	h.logger.Error("500 status ", err)
	w.WriteHeader(http.StatusInternalServerError)
}
