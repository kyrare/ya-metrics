package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kyrare/ya-metrics/internal/infrastructure/metrics"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestHandler_Home(t *testing.T) {
	logger := zap.New(nil)
	sugar := *logger.Sugar()
	storage := metrics.NewMemStorage("")

	ts := httptest.NewServer(ServerRouter(storage, false, sugar))
	defer ts.Close()

	type want struct {
		contentType string
		statusCode  int
	}

	tests := []struct {
		name    string
		request string
		method  string
		want    want
	}{
		{
			name:    "Get metric",
			request: "",
			method:  "GET",
			want: want{
				contentType: "text/html; charset=utf-8",
				statusCode:  http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, _ := testRequest(t, ts, tt.method, tt.request)

			defer func() {
				err := r.Body.Close()
				require.NoError(t, err)
			}()

			assert.Equal(t, tt.want.statusCode, r.StatusCode)
			assert.Equal(t, tt.want.contentType, r.Header.Get("Content-Type"))
		})
	}
}

func TestHandler_Get(t *testing.T) {
	logger := zap.New(nil)
	sugar := *logger.Sugar()
	storage := metrics.NewMemStorage("")

	ts := httptest.NewServer(ServerRouter(storage, false, sugar))
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
		want    want
	}{
		{
			name:    "Check not exists metric",
			request: "/value/gauge/foo",
			method:  "GET",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusNotFound,
				response:    "Not Found\n",
			},
		},
		{
			name:    "Check not exists metric type",
			request: "/value/foo/foo",
			method:  "GET",
			want: want{
				contentType: "text/plain; charset=utf-8",
				statusCode:  http.StatusBadRequest,
				response:    "Bad Request\n",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, b := testRequest(t, ts, tt.method, tt.request)

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

func TestHandler_Update(t *testing.T) {
	logger := zap.New(nil)
	sugar := *logger.Sugar()
	storage := metrics.NewMemStorage("")

	ts := httptest.NewServer(ServerRouter(storage, false, sugar))
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
			r, _ := testRequest(t, ts, "POST", tt.request)

			defer func() {
				err := r.Body.Close()
				require.NoError(t, err)
			}()

			assert.Equal(t, tt.want.statusCode, r.StatusCode)
			assert.Equal(t, tt.want.contentType, r.Header.Get("Content-Type"))
		})
	}
}

func testRequest(t *testing.T, ts *httptest.Server, method, path string) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, nil)
	require.NoError(t, err)

	r, err := ts.Client().Do(req)
	require.NoError(t, err)

	body, err := io.ReadAll(r.Body)
	require.NoError(t, err)

	return r, string(body)
}
