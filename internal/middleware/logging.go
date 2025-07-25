package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// responseWriter is a wrapper around http.ResponseWriter that captures the status code.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func LoggingMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rw := newResponseWriter(w)

			// Serve the request
			next.ServeHTTP(rw, r)

			// Log the request details after it's completed
			logger.Info("request completed",
				"method", r.Method,
				"uri", r.RequestURI,
				"ip", r.RemoteAddr,
				"user_agent", r.UserAgent(),
				"status", rw.statusCode,
				"duration", time.Since(start),
			)
		})
	}
}
