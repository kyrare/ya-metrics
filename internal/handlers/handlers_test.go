package handlers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUpdateHandle(t *testing.T) {
	ts := httptest.NewServer(ServerRouter())
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
				if err != nil {
					panic(err)
				}
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
