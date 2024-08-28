package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// WithLogging мидлвара для логирования запросов
func WithLogging(logger zap.SugaredLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		logFn := func(w http.ResponseWriter, r *http.Request) {
			// функция Now() возвращает текущее время
			start := time.Now()

			// эндпоинт /ping
			uri := r.RequestURI
			// метод запроса
			method := r.Method

			// точка, где выполняется хендлер pingHandler
			h.ServeHTTP(w, r) // обслуживание оригинального запроса

			// Since возвращает разницу во времени между start
			// и моментом вызова Since. Таким образом можно посчитать
			// время выполнения запроса.
			duration := time.Since(start)

			// отправляем сведения о запросе в zap
			logger.Infoln(
				"uri", uri,
				"method", method,
				"duration", duration,
			)
		}

		// возвращаем функционально расширенный хендлер
		return http.HandlerFunc(logFn)
	}
}
