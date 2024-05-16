package handlers

import (
	"net/http"
)

func (h *Handler) Ping(w http.ResponseWriter, r *http.Request) {
	if err := h.DB.Ping(); err == nil {
		w.WriteHeader(200)
	} else {
		h.logger.Error("500 status ", err)
		w.WriteHeader(500)
	}
}
