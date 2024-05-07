package middlewares

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (c compressWriter) Write(b []byte) (int, error) {
	return c.Writer.Write(b)
}

func Compress(h http.Handler) http.Handler {
	logFn := func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			h.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestSpeed)
		if err != nil {
			io.WriteString(w, err.Error())
			return
		}

		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		h.ServeHTTP(compressWriter{ResponseWriter: w, Writer: gz}, r)
	}

	// возвращаем функционально расширенный хендлер
	return http.HandlerFunc(logFn)
}
