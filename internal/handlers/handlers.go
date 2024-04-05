package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/kyrare/ya-metrics/internal/metrics"
	"github.com/kyrare/ya-metrics/internal/storage"
	"net/http"
	"strconv"
)

type Handler struct {
	storage storage.Storage
}

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
		panic(err)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	metricType := metrics.MetricType(chi.URLParam(r, "metricType"))
	metric := chi.URLParam(r, "metric")

	if metricType != metrics.TypeGauge && metricType != metrics.TypeCounter {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	value, ok := h.storage.Get(metricType, metric)

	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	_, err := w.Write([]byte(strconv.FormatFloat(value, 'f', -1, 64)))

	if err != nil {
		panic(err)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")

	metricType := metrics.MetricType(chi.URLParam(r, "metricType"))

	if metricType != metrics.TypeGauge && metricType != metrics.TypeCounter {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	metric := chi.URLParam(r, "metric")
	value, err := strconv.ParseFloat(chi.URLParam(r, "value"), 64)

	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if metricType == metrics.TypeGauge {
		h.storage.UpdateGauge(metric, value)
	}

	if metricType == metrics.TypeCounter {
		h.storage.UpdateCounter(metric, value)
	}
}

func NewHandler(storage storage.Storage) *Handler {
	return &Handler{storage: storage}
}

func ServerRouter() chi.Router {
	r := chi.NewRouter()

	s := storage.NewMemStorage()
	h := NewHandler(s)

	r.Get("/", h.Home)
	r.Get("/value/{metricType}/{metric}", h.Get)
	r.Post("/update/{metricType}/{metric}/{value}", h.Update)

	return r
}
