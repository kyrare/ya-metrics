package handlers

import (
	"net/http"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	err := h.DB.Ping()

	if err == nil {
		w.WriteHeader(200)
		return
	}

	h.logger.Error("500 status ", err)
	w.WriteHeader(500)
}
