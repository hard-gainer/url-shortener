package api

import (
	"log/slog"
	"net/http"
	"time"
)

// LoggingMiddleware loggs information about HTTP request
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := newLoggingResponseWriter(w)
		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		slog.Info("http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", lrw.statusCode,
			"duration_ms", duration.Milliseconds(),
		)
	})
}

// loggingResponseWriter creates a wrap of http.ResponseWriter to trace status code 
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// newLoggingResponseWriter creates new loggingResponseWriter
func newLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

// WriteHeader intercepts status code 
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
