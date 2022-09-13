package internalhttp

import (
	"net/http"
	"time"
)

type loggingWriter struct {
	http.ResponseWriter
	statusCode int
}

func newLoggingResponseWriter(w http.ResponseWriter) *loggingWriter {
	return &loggingWriter{w, http.StatusOK}
}

func (lrw *loggingWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (logMiddleware *loggingMiddleware) Process(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		logWriter := newLoggingResponseWriter(w)
		next.ServeHTTP(w, r)
		logMiddleware.logger.LogHTTPRequest(r, logWriter.statusCode, time.Since(start))
	})
}
