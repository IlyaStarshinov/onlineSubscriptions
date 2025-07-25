package handler

import (
	"log"
	"net/http"
	"time"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	status int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.status = code
	lrw.ResponseWriter.WriteHeader(code)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{w, http.StatusOK}

		log.Printf("[REQ]  %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)

		log.Printf("[RESP] %s %s completed in %v with %d", r.Method, r.URL.Path, duration, lrw.status)
	})
}
