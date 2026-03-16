package logger

import (
	"log/slog"
	"net/http"
)

func New(log *slog.Logger) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := r.Header.Get("X-Request-ID")

			log.Info(
				"request",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.String("request_id", reqID),
			)

			next.ServeHTTP(w, r)
		})
	}
}
