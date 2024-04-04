package handlers

import (
	"github.com/kyrare/ya-metrics/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandle(t *testing.T) {
	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name    string
		request string
		storage storage.Storage
		want    want
	}{
		{
			name:    "Add metric",
			request: "/update/gauge/foo/100",
			storage: &storage.MemStorage{
				Gauges:   make(map[string]float64),
				Counters: make(map[string][]float64),
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  200,
			},
		},
		{
			name:    "Call without type",
			request: "/update/foo/100",
			storage: &storage.MemStorage{
				Gauges:   make(map[string]float64),
				Counters: make(map[string][]float64),
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  404,
			},
		},
		{
			name:    "Call incorrect type",
			request: "/update/bar/foo/100",
			storage: &storage.MemStorage{
				Gauges:   make(map[string]float64),
				Counters: make(map[string][]float64),
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
		},
		{
			name:    "Call incorrect value",
			request: "/update/gauge/foo/bar",
			storage: &storage.MemStorage{
				Gauges:   make(map[string]float64),
				Counters: make(map[string][]float64),
			},
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  400,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			request := httptest.NewRequest(http.MethodPost, tt.request, nil)
			w := httptest.NewRecorder()

			UpdateHandle(w, request, tt.storage)

			result := w.Result()
			defer func() {
				err := result.Body.Close()
				if err != nil {
					panic(err)
				}
			}()

			assert.Equal(t, tt.want.statusCode, result.StatusCode)
			assert.Equal(t, tt.want.contentType, result.Header.Get("Content-Type"))
		})
	}
}
