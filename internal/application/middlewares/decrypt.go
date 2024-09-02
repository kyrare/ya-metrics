package middlewares

import (
	"bytes"
	"io"
	"net/http"

	"github.com/kyrare/ya-metrics/internal/infrastructure/encrypt"
	"go.uber.org/zap"
)

// Decrypt мидлвара для расшифровки запросов
func Decrypt(cryptoKey string, logger zap.SugaredLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			if cryptoKey == "" {
				h.ServeHTTP(w, r)
				return
			}

			logger.Infof("Передан путь к приватному файлу %s, начинаем дешифровку", cryptoKey)

			data, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			defer r.Body.Close()

			data, err = encrypt.Decrypt(data, cryptoKey)
			if err != nil {
				logger.Error(err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(data))

			h.ServeHTTP(w, r)
		}

		// возвращаем функционально расширенный хендлер
		return http.HandlerFunc(logFn)
	}
}
