package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

// Home получание главной страницы со списком всех метрик
func (h *Handler) Home(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)

	var result string

	result += "<h2>Gauges:</h2>"

	for metric, value := range h.storage.GetGauges() {
		result += fmt.Sprintf("<div><b>%s</b>: %s<div>", metric, strconv.FormatFloat(value, 'f', -1, 64))
	}

	result += "<h2>Counters:</h2>"

	for metric, value := range h.storage.GetCounters() {
		result += fmt.Sprintf("<div><b>%s</b>: %s<div>", metric, strconv.FormatFloat(value, 'f', -1, 64))
	}

	_, err := w.Write([]byte(result))
	if err != nil {
		h.logger.Error(err)
	}
}
