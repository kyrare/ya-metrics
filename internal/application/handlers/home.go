package handlers

import (
	"fmt"
	"net/http"
	"strconv"
)

func (h *Handler) Home(w http.ResponseWriter, _ *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")

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
		fmt.Println(err)
	}
}
