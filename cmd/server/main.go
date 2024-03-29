package main

import (
	"net/http"
	"strconv"
	"strings"
)

type Storage interface {
	UpdateGauge(metric string, value float64)
	UpdateCounter(metric string, value float64)
}

type MemStorage struct {
	Gauges   map[string]float64
	Counters map[string][]float64
}

func (storage *MemStorage) UpdateGauge(metric string, value float64) {
	if storage.Gauges == nil {
		storage.Gauges = make(map[string]float64)
	}

	storage.Gauges[metric] = value
}

func (storage *MemStorage) UpdateCounter(metric string, value float64) {
	if storage.Counters == nil {
		storage.Counters = make(map[string][]float64)
	}

	storage.Counters[metric] = append(storage.Counters[metric], value)
}

func main() {
	var storage MemStorage

	mux := http.NewServeMux()

	mux.Handle("/update/", postMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		updateHandle(w, r, &storage)
	})))

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}

func updateHandle(w http.ResponseWriter, r *http.Request, storage Storage) {
	parts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")

	if len(parts) != 4 {
		http.NotFound(w, r)
		return
	}

	metricType := parts[1]
	metric := parts[2]
	valueStr := parts[3]

	if metricType != "gauge" && metricType != "counter" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	value, err := strconv.ParseFloat(valueStr, 64)

	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	if metricType == "gauge" {
		storage.UpdateGauge(metric, value)
	}

	if metricType == "counter" {
		storage.UpdateCounter(metric, value)
	}
}

func postMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
