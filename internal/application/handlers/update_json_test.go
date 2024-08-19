package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyrare/ya-metrics/internal/infrastructure/connection"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHandler_UpdateJSON(t *testing.T) {
	logger := zap.New(nil)
	sugar := *logger.Sugar()
	storage, err := metrics.NewMemStorage("", sugar)
	assert.NoError(t, err, "Не удалось создать storage")

	db, err := connection.New("", sugar)
	assert.NoError(t, err, "Не удалось создать соединение с БД")

	ts := httptest.NewServer(ServerRouter(storage, db, false, false, "", sugar))
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
		response    string
	}

	tests := []struct {
		name    string
		request string
		method  string
		body    string
		want    want
	}{
		{
			name:    "Check incorrect json",
			request: "/update/",
			method:  "POST",
			body:    "",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "Bad Request\n",
			},
		},
		{
			name:    "Check add metric",
			request: "/update/",
			method:  "POST",
			body:    "{\"id\":\"TestAlloc\",\"type\":\"gauge\",\"value\":2.2}",
			want: want{
				contentType: "application/json",
				statusCode:  http.StatusOK,
				response:    "{\"id\":\"TestAlloc\",\"type\":\"gauge\",\"value\":2.2}",
			},
		},
		{
			name:    "Check not exists metric type",
			request: "/update/",
			method:  "POST",
			body:    "{\"id\":\"Alloc\",\"type\":\"testGauge\",\"value\":2.2}",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "Bad Request\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, b := testRequest(t, ts, tt.method, tt.request, bytes.NewBuffer([]byte(tt.body)))

			defer func() {
				err := r.Body.Close()
				require.NoError(t, err)
			}()

			assert.Equal(t, tt.want.statusCode, r.StatusCode)
			assert.Equal(t, tt.want.contentType, r.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.response, b)
		})
	}
}
