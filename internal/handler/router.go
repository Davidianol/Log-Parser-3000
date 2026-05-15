package handler

import (
	"log/slog"
	"net/http"
	"time"
)

func NewRouter(h *Handler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/parse/", h.Parse)
	mux.HandleFunc("/api/topology/", h.GetTopology)
	mux.HandleFunc("/api/node/", h.GetNode)
	mux.HandleFunc("/api/port/", h.GetPorts)
	mux.HandleFunc("/api/log/", h.GetLog)

	return loggingMiddleware(mux)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.status,
			"duration", time.Since(start),
		)
	})
}

// responseWriter -- перехватчик статусов кодов
type responseWriter struct {
	http.ResponseWriter
	status int
}

func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}
