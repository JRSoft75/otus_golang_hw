package internalhttp

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"time"
)

func LoggingMiddleware(logger Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now() // Засекаем время начала обработки запроса

			// Создаем объект ResponseWriter, чтобы перехватить код ответа
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			// Вызываем следующий обработчик в цепочке
			next.ServeHTTP(ww, r)

			// Вычисляем latency (время обработки запроса)
			latency := time.Since(start)

			ip := r.RemoteAddr
			dateTime := start.Format(time.RFC3339)
			method := r.Method
			path := r.URL.Path
			httpVersion := fmt.Sprintf("HTTP/%d.%d", r.ProtoMajor, r.ProtoMinor)
			statusCode := ww.Status()
			userAgent := r.UserAgent()

			logMessage := fmt.Sprintf(
				"IP: %s | Date: %s | Method: %s | Path: %s | HTTP Version: %s | Status Code: %d | Latency: %v | User-Agent: %s",
				ip, dateTime, method, path, httpVersion, statusCode, latency, userAgent,
			)

			logger.Info(logMessage)
		})
	}
}
