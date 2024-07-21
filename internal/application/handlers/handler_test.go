package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func testRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string) {
	req, err := http.NewRequest(method, ts.URL+path, body)
	require.NoError(t, err)

	r, err := ts.Client().Do(req)
	require.NoError(t, err)

	b, err := io.ReadAll(r.Body)
	require.NoError(t, err)

	return r, string(b)
}
