package middlewares

import (
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func (c compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}

func Decompress(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.Header.Get("Content-Encoding"))

		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		zr, err := gzip.NewReader(r.Body)

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		cr := &compressReader{
			r:  r.Body,
			zr: zr,
		}

		r.Body = cr
		defer cr.Close()

		h.ServeHTTP(w, r)
	}

	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}
