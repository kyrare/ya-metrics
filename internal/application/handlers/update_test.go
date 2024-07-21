package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyrare/ya-metrics/internal/infrastructure/connection"
	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHandler_Update(t *testing.T) {
	logger := zap.New(nil)
	sugar := *logger.Sugar()
	storage, err := metrics.NewMemStorage("", sugar)
	assert.NoError(t, err, "Не удалось создать storage")

	db, err := connection.New("", sugar)
	assert.NoError(t, err, "Не удалось создать соединение с БД")

	ts := httptest.NewServer(ServerRouter(storage, db, false, false, sugar))
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name    string
		request string
		want    want
	}{
		{
			name:    "Add metric",
			request: "/update/gauge/foo/100",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusOK,
			},
		},
		{
			name:    "Call without type",
			request: "/update/foo/100",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
			},
		},
		{
			name:    "Call incorrect type",
			request: "/update/bar/foo/100",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
		{
			name:    "Call incorrect value",
			request: "/update/gauge/foo/bar",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := testRequest(t, ts, "POST", tt.request, nil)

			defer func() {
				err := r.Body.Close()
				require.NoError(t, err)
			}()

			assert.Equal(t, tt.want.statusCode, r.StatusCode)
			assert.Equal(t, tt.want.contentType, r.Header.Get("Content-Type"))
		})
	}
}
